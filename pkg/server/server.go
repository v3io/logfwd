package server

import (
	"github.com/nuclio/logger"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/v3io/go-errors"
	"github.com/v3io/logfwd/pkg/record"
	"github.com/v3io/logfwd/pkg/rules"
	"github.com/valyala/fasthttp"
)

type Records []record.LogRecord

type Server struct {
	log logger.Logger

	configuration *rules.RuleConfig

	listenAddress string
}

func (s *Server) handleRecords(body []byte) error {
	decoder := record.NewArrayDecoder(s.log.GetChild("decoder"))
	records, err := decoder.FromByteArray(body)
	if err != nil {
		return errors.Wrap(err, "Unable to decode posted data")
	}
	for _, data := range records {
		if err := s.configuration.Send(&data); err != nil {
			s.log.WarnWith("Error when handling record",
				"error", err,
				"data", data)
		}
	}
	return nil
}

func (s *Server) handler(c *routing.Context) error {
	data := c.PostBody()
	return s.handleRecords(data)
}

func (s *Server) Run() error {
	router := routing.New()
	router.Any("/", s.handler)

	return fasthttp.ListenAndServe(s.listenAddress, router.HandleRequest)
}

func NewServer(log logger.Logger, listenAddress string, rules *rules.RuleConfig) (*Server, error) {
	s := &Server{
		log:           log.GetChild("server"),
		listenAddress: listenAddress,
		configuration: rules,
	}

	return s, nil
}
