package user_context

// <summary>
// Ctx
// <summary>
type Ctx interface {
	SetContext(any any)
	GetContext() any
	SetUserId(userId string)
	GetUserId() string
	SetPlatformId(platformId int)
	GetPlatformId() int
	SetToken(token string)
	GetToken() string
	SetSession(session string)
	GetSession() string
}

// <summary>
// context
// <summary>
type context struct {
	userId     string
	platformId int
	token      string
	session    string
	context    any
}

func NewCtx() Ctx {
	s := &context{}
	return s
}

func (s *context) SetContext(any any) {
	s.context = any
}

func (s *context) GetContext() any {
	return s.context
}

func (s *context) SetUserId(userId string) {
	s.userId = userId
}

func (s *context) GetUserId() string {
	return s.userId
}

func (s *context) SetPlatformId(platformId int) {
	s.platformId = platformId
}

func (s *context) GetPlatformId() int {
	return s.platformId
}

func (s *context) SetToken(token string) {
	s.token = token
}

func (s *context) GetToken() string {
	return s.token
}

func (s *context) SetSession(session string) {
	s.session = session
}

func (s *context) GetSession() string {
	return s.session
}
