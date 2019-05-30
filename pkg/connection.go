package pkg

type Connection struct {
	Host string
	User string
	Password string
	Port string
}

func NewConnection() *Connection {
	return &Connection{
		Host: "127.0.0.1",
		User: "root",
		Password: "",
		Port: "3306",
	}
}