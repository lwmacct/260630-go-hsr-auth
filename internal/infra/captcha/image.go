package captcha

import (
	"context"
	"crypto/sha256"
	"errors"
	"image/color"
	"strings"
	"sync"
	"time"

	"github.com/golang-module/base64Captcha/driver"

	"github.com/lwmacct/260630-go-hsr-auth/internal/infra/token"
	"github.com/lwmacct/260630-go-hsr-auth/internal/service"
)

const imageDefaultTTL = 2 * time.Minute

var errImageChallengeLimitExceeded = errors.New("challenge limit exceeded")

type ImageProvider struct {
	mu         sync.Mutex
	challenges map[string]imageChallenge
	driver     driver.Driver
	ttl        time.Duration
	maxItems   int
}

type imageChallenge struct {
	answerHash [32]byte
	expiresAt  time.Time
}

func NewImageProvider(maxItems int) *ImageProvider {
	return &ImageProvider{
		challenges: make(map[string]imageChallenge),
		driver: driver.NewDriverString(driver.DriverString{
			Width:           180,
			Height:          56,
			Length:          4,
			NoiseCount:      12,
			ShowLineOptions: driver.OptionShowHollowLine | driver.OptionShowSlimeLine,
			Source:          "23456789ABCDEFGHJKLMNPQRSTUVWXYZ",
			BgColor:         &color.RGBA{R: 248, G: 250, B: 252, A: 255},
		}),
		ttl:      imageDefaultTTL,
		maxItems: maxItems,
	}
}

func (p *ImageProvider) Name() string {
	return service.AuthChallengeProviderImage
}

func (p *ImageProvider) PublicConfig() service.AuthChallengePublicConfig {
	return service.AuthChallengePublicConfig{Provider: service.AuthChallengeProviderImage}
}

func (p *ImageProvider) Create(context.Context, service.AuthChallengeInput) (*service.AuthChallenge, error) {
	_, content, answer := p.driver.GenerateCaptcha()
	image, err := p.driver.DrawCaptcha(content)
	if err != nil {
		return nil, err
	}
	id, expiresAt, err := p.put(answer)
	if err != nil {
		if errors.Is(err, errImageChallengeLimitExceeded) {
			return nil, service.ErrAuthChallengeLimitExceeded
		}
		return nil, err
	}
	return &service.AuthChallenge{
		Provider:    service.AuthChallengeProviderImage,
		ChallengeID: id,
		Image:       image.Encoder(),
		ExpiresAt:   expiresAt,
	}, nil
}

func (p *ImageProvider) Verify(_ context.Context, response service.AuthChallengeAnswer, _ service.AuthChallengeInput) error {
	if !p.verifyAndDelete(response.ChallengeID, response.Answer) {
		return service.ErrAuthChallengeInvalid
	}
	return nil
}

func (p *ImageProvider) put(answer string) (string, time.Time, error) {
	id, err := token.NewWithPrefix("cap")
	if err != nil {
		return "", time.Time{}, err
	}
	now := time.Now().UTC()
	expiresAt := now.Add(p.ttl)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cleanupLocked(now)
	if p.maxItems > 0 && len(p.challenges) >= p.maxItems {
		return "", time.Time{}, errImageChallengeLimitExceeded
	}
	p.challenges[id] = imageChallenge{answerHash: imageAnswerHash(answer), expiresAt: expiresAt}
	return id, expiresAt, nil
}

func (p *ImageProvider) verifyAndDelete(id string, answer string) bool {
	id = strings.TrimSpace(id)
	if id == "" || strings.TrimSpace(answer) == "" {
		return false
	}
	now := time.Now().UTC()
	p.mu.Lock()
	defer p.mu.Unlock()
	challenge, ok := p.challenges[id]
	if !ok {
		return false
	}
	delete(p.challenges, id)
	if !challenge.expiresAt.After(now) {
		return false
	}
	return challenge.answerHash == imageAnswerHash(answer)
}

func (p *ImageProvider) cleanupLocked(now time.Time) {
	for id, challenge := range p.challenges {
		if !challenge.expiresAt.After(now) {
			delete(p.challenges, id)
		}
	}
}

func imageAnswerHash(answer string) [32]byte {
	return sha256.Sum256([]byte(strings.ToUpper(strings.TrimSpace(answer))))
}
