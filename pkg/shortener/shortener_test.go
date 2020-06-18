package shortener

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mock_shortener "github.com/zikoel/shortener/mocks"
)

func TestURLFromKey_URLNotValid(t *testing.T) {

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	const url string = "google.com"

	f := mock_shortener.NewMockFinder(ctrl)
	p := mock_shortener.NewMockPersister(ctrl)

	f.EXPECT().Lookup(nil).Times(0)
	p.EXPECT().Persist(nil, nil).Times(0)

	s, err := CreateShortener(f, p, 1234)

	if err != nil {
		t.Error("Error on CreateShortener")
	}

	_, err = s.KeyFromURL(url, "")

	if err == nil {
		t.Error("We expect an error")
	}
}

func TestURLFromKey_keyNotSuggestedNoCollision(t *testing.T) {

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	const url string = "http://www.google.com"

	f := mock_shortener.NewMockFinder(ctrl)
	p := mock_shortener.NewMockPersister(ctrl)

	f.EXPECT().Lookup( gomock.Any() ).Return("", errors.New("err"))
	p.EXPECT().Persist(gomock.Any(), gomock.Eq(url)).Return(nil).Times(1)
	p.EXPECT().Persist(gomock.Any(), gomock.Eq(0)).Return(nil).Times(1)

	s, err := CreateShortener(f, p, 1234)

	if err != nil {
		t.Error("Error on CreateShortener")
	}

	key, err := s.KeyFromURL(url, "")

	if err != nil {
		t.Error("Error on KeyFromURL")
	}

	if len(key) != 8 {
		t.Error("key length not correct")
	}
}

func TestURLFromKey_keySuggestedCollision(t *testing.T) {

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	const url string = "http://www.google.com"

	f := mock_shortener.NewMockFinder(ctrl)
	p := mock_shortener.NewMockPersister(ctrl)

	f.EXPECT().Lookup( gomock.Eq("foo") ).Return(url, nil)
	p.EXPECT().Persist(nil, nil).Times(0)

	s, err := CreateShortener(f, p, 1234)

	if err != nil {
		t.Error("Error on CreateShortener")
	}

	key, err := s.KeyFromURL(url, "foo")

	if err != nil {
		t.Error("Error on KeyFromURL")
	}

	if key != "foo" {
		t.Error("unexpected key")
	}
}

func TestURLFromKey_keySuggestedNoCollision(t *testing.T) {

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	const url string = "http://www.google.com"
	const keySuggested = "foo"

	f := mock_shortener.NewMockFinder(ctrl)
	p := mock_shortener.NewMockPersister(ctrl)

	f.EXPECT().Lookup( gomock.Eq(keySuggested) ).Return("", errors.New("err"))
	p.EXPECT().Persist(gomock.Eq(keySuggested), gomock.Eq(url)).Return(nil).Times(1)
	p.EXPECT().Persist(gomock.Eq(keySuggested+":count"), gomock.Eq(0)).Return(nil).Times(1)

	s, err := CreateShortener(f, p, 1234)

	if err != nil {
		t.Error("Error on CreateShortener")
	}

	key, err := s.KeyFromURL(url, "foo")

	if err != nil {
		t.Error("Error on KeyFromURL")
	}

	if key != "foo" {
		t.Error("unexpected key")
	}
}

func TestURLFromKey_keyNoSuggestedCollision(t *testing.T) {

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	const url string = "http://www.google.com"
	const anotherURL string = "http://www.hackernews.com"
	const anotherURL2 string = "http://www.medium.com"

	f := mock_shortener.NewMockFinder(ctrl)
	p := mock_shortener.NewMockPersister(ctrl)

	f.EXPECT().Lookup( gomock.Any() ).Return(anotherURL, nil).Times(1) // Simulate first collision
	f.EXPECT().Lookup( gomock.Any() ).Return(anotherURL, nil).Times(1) // Simulate second collision
	f.EXPECT().Lookup( gomock.Any() ).Return("", errors.New("err")).Times(1) // The third finally goes
	p.EXPECT().Persist(gomock.Any(), gomock.Eq(url)).Return(nil).Times(1)
	p.EXPECT().Persist(gomock.Any(), gomock.Eq(0)).Return(nil).Times(1)

	s, err := CreateShortener(f, p, 1234)

	if err != nil {
		t.Error("Error on CreateShortener")
	}

	key, err := s.KeyFromURL(url, "")

	if err != nil {
		t.Error("Error on KeyFromURL")
	}

	if len(key) != 8 {
		t.Error("unexpected key length")
	}
}

func TestDeleteURLByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const key = "foobar"

	f := mock_shortener.NewMockFinder(ctrl)
	p := mock_shortener.NewMockPersister(ctrl)

	f.EXPECT().Lookup(nil).Times(0)
	p.EXPECT().Delete(gomock.Eq(key)).Return(nil).Times(1)
	p.EXPECT().Delete(gomock.Eq(key+":count")).Return(nil).Times(1)

	s, err := CreateShortener(f, p, 1234)

	if err != nil {
		t.Error("Error on CreateShortener")
	}

	err = s.DeleteURLByKey(key)

	if err != nil {
		t.Error("We expect any error here")
	}
}