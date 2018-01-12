package orders

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dynamicgo/config"
	"github.com/dynamicgo/slf4go"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/inwecrypto/ethdb"
)

// APIServer .
type APIServer struct {
	engine *gin.Engine
	slf4go.Logger
	laddr   string
	db      *xorm.Engine
	watcher *txWatcher
}

// Order .
type Order struct {
	ID          int64   `json:"-" xorm:"pk autoincr"`
	TX          string  `json:"tx" xorm:"notnull"`
	From        string  `json:"from" xorm:"index(from_to)"`
	To          string  `json:"to" xorm:"index(from_to)"`
	Asset       string  `json:"asset" xorm:"notnull"`
	Value       string  `json:"value" xorm:"notnull"`
	Blocks      uint64  `json:"blocks" xorm:""`
	CreateTime  *string `json:"createTime,omitempty" xorm:"TIMESTAMP notnull created"`
	ConfirmTime *string `json:"confirmTime" xorm:"TIMESTAMP"`
	Context     *string `json:"context,omitempty" xorm:"json"`
}

// NewAPIServer .
func NewAPIServer(conf *config.Config) (*APIServer, error) {

	db, err := initXORM(conf)

	if err != nil {
		return nil, err
	}

	if !conf.GetBool("orders.debug", true) {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	if conf.GetBool("orders.debug", true) {
		engine.Use(gin.Logger())
	}

	watcher, err := newTxWatcher(conf, db)

	if err != nil {
		return nil, err
	}

	server := &APIServer{
		engine:  engine,
		Logger:  slf4go.Get("eth-orders"),
		laddr:   conf.GetString("orders.laddr", ":8000"),
		db:      db,
		watcher: watcher,
	}

	server.makeRouters()

	return server, nil
}

func initXORM(conf *config.Config) (*xorm.Engine, error) {
	username := conf.GetString("orders.ethdb.username", "xxx")
	password := conf.GetString("orders.ethdb.password", "xxx")
	port := conf.GetString("orders.ethdb.port", "6543")
	host := conf.GetString("orders.ethdb.host", "localhost")
	scheme := conf.GetString("orders.ethdb.schema", "postgres")

	return xorm.NewEngine(
		"postgres",
		fmt.Sprintf(
			"user=%v password=%v host=%v dbname=%v port=%v sslmode=disable",
			username, password, host, scheme, port,
		),
	)
}

func (server *APIServer) makeRouters() {
	server.engine.POST("/wallet/:userid/:address", func(ctx *gin.Context) {
		if err := server.createWallet(ctx.Param("address"), ctx.Param("userid")); err != nil {
			server.ErrorF("create wallet error :%s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	})

	server.engine.DELETE("/wallet/:userid/:address", func(ctx *gin.Context) {
		if err := server.deleteWallet(ctx.Param("address"), ctx.Param("userid")); err != nil {
			server.ErrorF("create wallet error :%s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	})

	server.engine.POST("/order", func(ctx *gin.Context) {

		var order *Order

		if err := ctx.ShouldBindJSON(&order); err != nil {
			server.ErrorF("parse order error :%s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := server.createOrder(order); err != nil {
			server.ErrorF("create order error :%s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	})

	server.engine.GET("/order/:tx", func(ctx *gin.Context) {
		if orders, err := server.getOrder(ctx.Param("tx")); err != nil {
			server.ErrorF("get orders error :%s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, orders)
		}
	})

	server.engine.GET("/orders/:address/:asset/:offset/:size", func(ctx *gin.Context) {
		offset, err := parseInt(ctx, "offset")

		if err != nil {
			server.ErrorF("parse page parameter error :%s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		size, err := parseInt(ctx, "size")

		if err != nil {
			server.ErrorF("parse page parameter error :%s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		orders, err := server.getPagedOrders(ctx.Param("address"), ctx.Param("asset"), offset, size)

		if err != nil {
			server.ErrorF("get paged orders error :%s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, orders)

	})
}

func parseInt(ctx *gin.Context, name string) (int, error) {
	result, err := strconv.ParseInt(ctx.Param(name), 10, 32)

	return int(result), err
}

func (server *APIServer) createWallet(userid string, address string) error {

	wallet := &ethdb.TableWallet{
		Address: address,
		UserID:  userid,
	}

	_, err := server.db.Insert(wallet)
	return err
}

func (server *APIServer) deleteWallet(userid string, address string) error {

	wallet := &ethdb.TableWallet{
		Address: address,
		UserID:  userid,
	}

	_, err := server.db.Delete(wallet)
	return err
}

func (server *APIServer) createOrder(order *Order) error {

	tOrder := &ethdb.TableOrder{
		ID:      order.ID,
		TX:      order.TX,
		From:    order.From,
		To:      order.To,
		Asset:   order.Asset,
		Value:   order.Value,
		Context: order.Context,
		Blocks:  int64(order.Blocks),
	}

	_, err := server.db.Insert(tOrder)
	return err
}

func (server *APIServer) getOrder(tx string) ([]*Order, error) {

	server.DebugF("get order by tx %s", tx)

	torders := make([]*ethdb.TableOrder, 0)

	err := server.db.Where("t_x = ?", tx).Find(&torders)

	if err != nil {
		return make([]*Order, 0), err
	}

	orders := make([]*Order, 0)

	for _, torder := range torders {
		createTime := torder.CreateTime.Format(time.RFC3339Nano)

		var confirmTime *string

		if torder.ConfirmTime != nil {
			timestr := torder.ConfirmTime.Format(time.RFC3339Nano)
			confirmTime = &timestr
		}

		orders = append(orders, &Order{
			TX:          torder.TX,
			From:        torder.From,
			To:          torder.To,
			Asset:       torder.Asset,
			Value:       torder.Value,
			Context:     torder.Context,
			CreateTime:  &createTime,
			ConfirmTime: confirmTime,
			Blocks:      uint64(torder.Blocks),
		})
	}

	return orders, nil
}

func (server *APIServer) getPagedOrders(address, asset string, offset, size int) ([]*Order, error) {

	server.DebugF("get address(%s) orders(%s) (%d,%d)", address, asset, offset, size)

	torders := make([]*ethdb.TableOrder, 0)

	err := server.db.
		Where(`("from" = ? or "to" = ?) and asset = ?`, address, address, asset).
		Desc("create_time").
		Limit(size, offset).
		Find(&torders)

	if err != nil {
		return make([]*Order, 0), err
	}

	orders := make([]*Order, 0)

	for _, torder := range torders {
		createTime := torder.CreateTime.Format(time.RFC3339Nano)

		var confirmTime *string

		if torder.ConfirmTime != nil {
			timestr := torder.ConfirmTime.Format(time.RFC3339Nano)
			confirmTime = &timestr
		}

		orders = append(orders, &Order{
			TX:          torder.TX,
			From:        torder.From,
			To:          torder.To,
			Asset:       torder.Asset,
			Value:       torder.Value,
			Context:     torder.Context,
			CreateTime:  &createTime,
			ConfirmTime: confirmTime,
			Blocks:      uint64(torder.Blocks),
		})
	}

	return orders, nil
}

// Run run http service
func (server *APIServer) Run() error {
	go server.watcher.Run()
	return server.engine.Run(server.laddr)
}
