package health

import "github.com/charmbracelet/lipgloss"

var (
	// Colors (design tokens)
	colorPrimaryBorder = lipgloss.Color("#6366F1")
	colorTitle         = lipgloss.Color("#A5B4FC")
	colorWhite         = lipgloss.Color("#FFFFFF")
	colorPass          = lipgloss.Color("#22C55E")
	colorFail          = lipgloss.Color("#EF4444")
	colorWarn          = lipgloss.Color("#F59E0B")
	colorWarnSoft      = lipgloss.Color("#FCD34D")
	colorCritSoft      = lipgloss.Color("#FCA5A5")
	colorMuted         = lipgloss.Color("#D1D5DB")

	// Layout tokens
	cardWidth      = 58
	cardPaddingLR  = 2
	badgePaddingLR = 1
	labelPaddingL  = 1
)

var (
	headerCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimaryBorder).
			Padding(0, cardPaddingLR).
			Width(cardWidth)

	headerTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorTitle)

	headerSubTitleStyle = lipgloss.NewStyle().
				Foreground(colorMuted)

	domainStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite)

	// Status row badges
	passBadge = lipgloss.NewStyle().
			Background(colorPass).
			Foreground(colorWhite).
			Padding(0, badgePaddingLR).
			Bold(true)

	failBadge = lipgloss.NewStyle().
			Background(colorFail).
			Foreground(colorWhite).
			Padding(0, badgePaddingLR).
			Bold(true)

	checkLabel = lipgloss.NewStyle().
			Foreground(colorMuted).
			PaddingLeft(labelPaddingL)

	// Section titles
	warnTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWarn)

	critTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorFail)

	okStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPass)

	// Bullet items
	warnBullet = lipgloss.NewStyle().Foreground(colorWarnSoft)
	critBullet = lipgloss.NewStyle().Foreground(colorCritSoft)

	// Report extras
	dividerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#374151"))
	scoreBarFill  = lipgloss.NewStyle().Foreground(colorPass)
	scoreBarEmpty = lipgloss.NewStyle().Foreground(lipgloss.Color("#374151"))
	scoreBarLabel = lipgloss.NewStyle().Foreground(colorMuted)
	hintStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Italic(true)
)
