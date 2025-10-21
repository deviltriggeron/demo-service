package entity

type Order struct {
	OrderUID        string `json:"order_uid" validate:"required"`
	TrackNumber     string `json:"track_number" validate:"required"`
	Entry           string `json:"entry" validate:"required"`
	Locale          string `json:"locale" validate:"required"`
	InternalSig     string `json:"internal_signature"`
	CustomerID      string `json:"customer_id" validate:"required"`
	DeliveryService string `json:"delivery_service" validate:"required"`
	Shardkey        string `json:"shardkey"`
	SmID            int    `json:"sm_id"`
	DateCreated     string `json:"date_created" validate:"required"`
	OofShard        string `json:"oof_shard"`

	Delivery Delivery `json:"delivery" validate:"required,dive"`
	Payment  Payment  `json:"payment" validate:"required,dive"`
	Items    []Item   `json:"items" validate:"required,dive"`
}

type Delivery struct {
	OrderUID string `json:"order_uid"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Zip      string `json:"zip"`
	City     string `json:"city" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Region   string `json:"region"`
	Email    string `json:"email" validate:"required,email"`
}

type Payment struct {
	OrderUID     string `json:"order_uid"`
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required,len=3"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int    `json:"amount" validate:"gt=0"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	OrderUID    string `json:"order_uid"`
	ChrtID      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"gt=0"`
	Rid         string `json:"rid"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type KafkaMsg struct {
	Order  Order  `json:"order"`
	Method string `json:"method"`
}

type ConfigDB struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     string
	ServerPort       string
}

type ConfigBroker struct {
	Broker  string
	GroupID string
}
