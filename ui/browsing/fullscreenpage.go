package browsing

import (
	"github.com/dweymouth/supersonic/backend"
	"github.com/dweymouth/supersonic/backend/mediaprovider"
	"github.com/dweymouth/supersonic/sharedutil"
	"github.com/dweymouth/supersonic/ui/controller"
	"github.com/dweymouth/supersonic/ui/layouts"
	"github.com/dweymouth/supersonic/ui/util"
	"github.com/dweymouth/supersonic/ui/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FullscreenPage struct {
	widget.BaseWidget

	fullscreenPageState

	queue []*mediaprovider.Track

	card         *widgets.LargeNowPlayingCard
	nowPlayingID string
}

type fullscreenPageState struct {
	contr   *controller.Controller
	pool    *util.WidgetPool
	pm      *backend.PlaybackManager
	im      *backend.ImageManager
	canRate bool
}

func NewFullscreenPage(
	contr *controller.Controller,
	pool *util.WidgetPool,
	im *backend.ImageManager,
	pm *backend.PlaybackManager,
	canRate bool,
) *FullscreenPage {
	a := &FullscreenPage{fullscreenPageState: fullscreenPageState{
		contr: contr, pool: pool, pm: pm, canRate: canRate,
	}}
	a.ExtendBaseWidget(a)

	a.card = widgets.NewLargeNowPlayingCard()

	a.Reload()
	return a
}

func (a *FullscreenPage) CreateRenderer() fyne.WidgetRenderer {
	container := container.NewGridWithColumns(2,
		container.New(&layouts.PercentPadLayout{
			LeftRightObjectPercent: .8,
			TopBottomObjectPercent: .8,
		}, a.card),
		layout.NewSpacer(),
	)
	return widget.NewSimpleRenderer(container)
}

func (a *FullscreenPage) Save() SavedPage {
	nps := a.fullscreenPageState
	return &nps
}

func (a *FullscreenPage) Route() controller.Route {
	return controller.NowPlayingRoute("")
}

var _ CanShowNowPlaying = (*FullscreenPage)(nil)

func (a *FullscreenPage) OnSongChange(song, lastScrobbledIfAny *mediaprovider.Track) {
	a.nowPlayingID = sharedutil.TrackIDOrEmptyStr(song)
	a.card.Update(song.Name, song.ArtistNames, song.ArtistIDs, song.Album)
}

func (a *FullscreenPage) Reload() {
	a.queue = a.pm.GetPlayQueue()
}

func (s *fullscreenPageState) Restore() Page {
	return NewFullscreenPage(s.contr, s.pool, s.im, s.pm, s.canRate)
}
