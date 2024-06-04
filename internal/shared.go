package shared

type Config struct {
	TraceServiceName  string
	GrpcPort          int
	PostgresHost      string
	PostgresPort      int
	PostgresUser      string
	PostgresPassword  string
	PostgresDb        string
	PostgresSchema    string
	DbMaxIdleConns    int
	DbMaxOpenConns    int
	KafkaVersion      string
	KafkaBroker       string
	RedisHost         string
	RedisPort         int
	RedisPassword     string
	RedisDb           int
	TemporalEndpoint  string
	TemporalNamespace string
}

// kafka topics (naming rules: "giaek-{consumer}-{action}")
const KafkaTopicSomething = "giaek-order-something"

// kafka group ids
const KafkaGroupIdOrder = "giaek-order"

// kafka message keys
const KafkaMsgKeyData = "data"

// temporal task queues
const TemporalTaskQueueOrder = "ORDER_TASK_QUEUE"
