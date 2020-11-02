package gservice

type service struct {
	host     string
	protocol RpcType
	factory  ConnFactory
	testFun  ConnTestFunc
}

func (s *service) fromTags(tags []string) {
	for _, tag := range tags {
		switch RpcType(tag) {
		case RpcTypeGRPC:
			s.protocol = RpcTypeGRPC
		}
	}
	if s.protocol == "" {
		s.protocol = RpcTypeGRPC
	}

	switch s.protocol {
	case RpcTypeGRPC:
		s.factory = GRpcFactory()
		s.testFun = GRpcConnTest()
	}
}
