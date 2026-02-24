package jpeg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	stdjpeg "image/jpeg"
	"math"
	"strconv"
	"strings"

	"github.com/user/gin-microservice-boilerplate/models"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	ErrInvalidChartPayload = errors.New("invalid chart payload")
)

const (
	defaultWidth   = 1280
	defaultHeight  = 720
	defaultQuality = 90
	minDimension   = 320
	maxDimension   = 4096
)

type chartStrategy struct{}

type chartPayload struct {
	Categories []string      `json:"categories"`
	Series     []chartSeries `json:"series"`
	Width      int           `json:"width,omitempty"`
	Height     int           `json:"height,omitempty"`
	Quality    int           `json:"quality,omitempty"`
}

type chartSeries struct {
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
}

func NewChartStrategy() *chartStrategy {
	return &chartStrategy{}
}

func (s *chartStrategy) Generate(ctx context.Context, req models.ExportRequest) (*models.ExportFile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	payload, err := decodeChartPayload(req.Payload)
	if err != nil {
		return nil, err
	}

	width := normalizeDimension(payload.Width, defaultWidth)
	height := normalizeDimension(payload.Height, defaultHeight)
	quality := normalizeQuality(payload.Quality)

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	drawBackground(canvas, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	drawChart(canvas, payload)

	var buffer bytes.Buffer
	if err := stdjpeg.Encode(&buffer, canvas, &stdjpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}

	return &models.ExportFile{
		FileName:    ensureJPEGFileName(req.FileName, "chart_export.jpeg"),
		ContentType: "image/jpeg",
		Data:        buffer.Bytes(),
	}, nil
}

func decodeChartPayload(raw interface{}) (*chartPayload, error) {
	var encoded []byte

	switch v := raw.(type) {
	case nil:
		return nil, fmt.Errorf("%w: payload is required", ErrInvalidChartPayload)
	case json.RawMessage:
		encoded = v
	case []byte:
		encoded = v
	default:
		var err error
		encoded, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("%w: malformed json", ErrInvalidChartPayload)
		}
	}

	var payload chartPayload
	if err := json.Unmarshal(encoded, &payload); err != nil {
		return nil, fmt.Errorf("%w: malformed json", ErrInvalidChartPayload)
	}

	if len(payload.Categories) == 0 {
		return nil, fmt.Errorf("%w: categories is required", ErrInvalidChartPayload)
	}
	if len(payload.Series) == 0 {
		return nil, fmt.Errorf("%w: series is required", ErrInvalidChartPayload)
	}

	for idx, series := range payload.Series {
		if strings.TrimSpace(series.Name) == "" {
			return nil, fmt.Errorf("%w: series[%d].name is required", ErrInvalidChartPayload, idx)
		}
		if len(series.Values) != len(payload.Categories) {
			return nil, fmt.Errorf(
				"%w: series[%d].values length must match categories length",
				ErrInvalidChartPayload,
				idx,
			)
		}
	}

	return &payload, nil
}

func drawChart(img *image.RGBA, payload *chartPayload) {
	plot := chartArea(img.Bounds())
	if plot.Dx() <= 0 || plot.Dy() <= 0 {
		return
	}

	drawGrid(img, plot, len(payload.Categories))
	drawAxes(img, plot)

	set := classifySeries(payload.Series)
	minValue, maxValue := minMaxValues(payload.Series)
	if almostEqual(minValue, maxValue) {
		maxValue = minValue + 1
	}

	drawYAxisLabels(img, plot, minValue, maxValue)
	drawXAxisLabels(img, plot, payload.Categories)

	if set.expectedMin != nil && set.expectedMax != nil {
		top := buildSeriesPoints(set.expectedMax.Values, len(payload.Categories), plot, minValue, maxValue)
		bottom := buildSeriesPoints(set.expectedMin.Values, len(payload.Categories), plot, minValue, maxValue)
		drawBand(img, top, bottom, color.RGBA{R: 148, G: 180, B: 234, A: 72})
		drawPolyline(img, top, color.RGBA{R: 128, G: 162, B: 220, A: 180})
		drawPolyline(img, bottom, color.RGBA{R: 128, G: 162, B: 220, A: 180})
	}

	colors := []color.RGBA{
		{R: 29, G: 78, B: 216, A: 255},
		{R: 16, G: 185, B: 129, A: 255},
		{R: 220, G: 38, B: 38, A: 255},
		{R: 168, G: 85, B: 247, A: 255},
	}

	if set.observed != nil {
		points := buildSeriesPoints(set.observed.Values, len(payload.Categories), plot, minValue, maxValue)
		drawPolyline(img, points, color.RGBA{R: 18, G: 138, B: 237, A: 255})
		drawObservedPoints(img, points, set)
	}

	for seriesIdx, series := range set.other {
		c := colors[seriesIdx%len(colors)]
		points := buildSeriesPoints(series.Values, len(payload.Categories), plot, minValue, maxValue)
		drawPolyline(img, points, c)
		for _, pt := range points {
			drawCircle(img, pt, 2, c)
		}
	}

	drawLegend(img, plot)
}

type seriesSet struct {
	observed    *chartSeries
	expectedMin *chartSeries
	expectedMax *chartSeries
	other       []*chartSeries
}

func classifySeries(series []chartSeries) seriesSet {
	set := seriesSet{}

	for i := range series {
		name := normalizeSeriesName(series[i].Name)
		switch {
		case strings.Contains(name, "expected") && (strings.Contains(name, "min") || strings.Contains(name, "lower")):
			set.expectedMin = &series[i]
		case strings.Contains(name, "expected") && (strings.Contains(name, "max") || strings.Contains(name, "upper")):
			set.expectedMax = &series[i]
		case strings.Contains(name, "observed") || strings.Contains(name, "download"):
			set.observed = &series[i]
		default:
			set.other = append(set.other, &series[i])
		}
	}

	if set.observed == nil {
		for i := range series {
			if &series[i] != set.expectedMin && &series[i] != set.expectedMax {
				set.observed = &series[i]
				break
			}
		}
	}

	return set
}

func normalizeSeriesName(name string) string {
	n := strings.TrimSpace(strings.ToLower(name))
	n = strings.ReplaceAll(n, " ", "_")
	n = strings.ReplaceAll(n, "-", "_")
	return n
}

func chartArea(bounds image.Rectangle) image.Rectangle {
	const (
		paddingLeft   = 96
		paddingRight  = 40
		paddingTop    = 40
		paddingBottom = 140
	)

	return image.Rect(
		bounds.Min.X+paddingLeft,
		bounds.Min.Y+paddingTop,
		bounds.Max.X-paddingRight,
		bounds.Max.Y-paddingBottom,
	)
}

func drawYAxisLabels(img *image.RGBA, plot image.Rectangle, minValue, maxValue float64) {
	const ticks = 5
	labelColor := color.RGBA{R: 82, G: 98, B: 116, A: 255}

	for i := 0; i <= ticks; i++ {
		value := maxValue - (maxValue-minValue)*float64(i)/float64(ticks)
		y := plot.Min.Y + int(float64(plot.Dy()-1)*float64(i)/float64(ticks))
		label := formatAxisValue(value)
		w := textWidth(label)
		drawText(img, plot.Min.X-12-w, y+4, label, labelColor)
	}
}

func drawXAxisLabels(img *image.RGBA, plot image.Rectangle, categories []string) {
	labelColor := color.RGBA{R: 82, G: 98, B: 116, A: 255}
	if len(categories) == 0 {
		return
	}

	maxTicks := 9
	step := 1
	if len(categories) > maxTicks {
		step = int(math.Ceil(float64(len(categories)) / float64(maxTicks)))
	}

	y := plot.Max.Y + 22
	lastDrawn := -1
	for i := 0; i < len(categories); i += step {
		label := shortCategoryLabel(categories[i])
		x := pointX(i, len(categories), plot)
		w := textWidth(label)
		drawText(img, x-w/2, y, label, labelColor)
		lastDrawn = i
	}

	if lastDrawn != len(categories)-1 {
		label := shortCategoryLabel(categories[len(categories)-1])
		x := pointX(len(categories)-1, len(categories), plot)
		w := textWidth(label)
		drawText(img, x-w/2, y, label, labelColor)
	}
}

func drawLegend(img *image.RGBA, plot image.Rectangle) {
	type legendItem struct {
		label string
		kind  string
	}

	items := []legendItem{
		{label: "Downloads", kind: "downloads"},
		{label: "Expected Range", kind: "range"},
		{label: "High", kind: "high"},
		{label: "Medium", kind: "medium"},
		{label: "Low", kind: "low"},
	}

	const (
		symbolWidth = 16
		itemGap     = 22
		textGap     = 8
	)

	totalWidth := 0
	widths := make([]int, len(items))
	for i, item := range items {
		w := symbolWidth + textGap + textWidth(item.label)
		widths[i] = w
		totalWidth += w
		if i < len(items)-1 {
			totalWidth += itemGap
		}
	}

	startX := plot.Min.X + (plot.Dx()-totalWidth)/2
	y := plot.Max.Y + 80
	textColor := color.RGBA{R: 82, G: 98, B: 116, A: 255}

	x := startX
	for i, item := range items {
		switch item.kind {
		case "downloads":
			blue := color.RGBA{R: 18, G: 138, B: 237, A: 255}
			drawLine(img, image.Pt(x, y-5), image.Pt(x+symbolWidth, y-5), blue)
			drawCircle(img, image.Pt(x+symbolWidth/2, y-5), 3, blue)
		case "range":
			fillRect(img, image.Rect(x, y-11, x+symbolWidth, y-3), color.RGBA{R: 148, G: 180, B: 234, A: 120})
		case "high":
			drawCircle(img, image.Pt(x+symbolWidth/2, y-7), 5, color.RGBA{R: 220, G: 38, B: 38, A: 255})
		case "medium":
			drawCircle(img, image.Pt(x+symbolWidth/2, y-7), 5, color.RGBA{R: 234, G: 88, B: 12, A: 255})
		case "low":
			drawCircle(img, image.Pt(x+symbolWidth/2, y-7), 5, color.RGBA{R: 217, G: 119, B: 6, A: 255})
		}

		drawText(img, x+symbolWidth+textGap, y, item.label, textColor)
		x += widths[i] + itemGap
	}
}

func drawBackground(img *image.RGBA, c color.Color) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, c)
		}
	}
}

func drawGrid(img *image.RGBA, plot image.Rectangle, categoryCount int) {
	gridColor := color.RGBA{R: 228, G: 232, B: 240, A: 255}
	const lines = 5

	maxX := plot.Max.X - 1
	maxY := plot.Max.Y - 1
	height := plot.Dy() - 1

	for i := 0; i <= lines; i++ {
		y := plot.Min.Y + int(float64(height)*float64(i)/float64(lines))
		drawLine(img, image.Pt(plot.Min.X, y), image.Pt(maxX, y), gridColor)
	}

	if categoryCount > 1 {
		for i := 0; i < categoryCount; i++ {
			x := pointX(i, categoryCount, plot)
			drawLine(img, image.Pt(x, plot.Min.Y), image.Pt(x, maxY), gridColor)
		}
	}
}

func drawAxes(img *image.RGBA, plot image.Rectangle) {
	axisColor := color.RGBA{R: 100, G: 116, B: 139, A: 255}
	maxX := plot.Max.X - 1
	maxY := plot.Max.Y - 1
	drawLine(img, image.Pt(plot.Min.X, maxY), image.Pt(maxX, maxY), axisColor)
	drawLine(img, image.Pt(plot.Min.X, plot.Min.Y), image.Pt(plot.Min.X, maxY), axisColor)
}

func minMaxValues(series []chartSeries) (float64, float64) {
	minValue := math.MaxFloat64
	maxValue := -math.MaxFloat64

	for _, s := range series {
		for _, value := range s.Values {
			if value < minValue {
				minValue = value
			}
			if value > maxValue {
				maxValue = value
			}
		}
	}
	return minValue, maxValue
}

func buildSeriesPoints(values []float64, total int, plot image.Rectangle, minValue, maxValue float64) []image.Point {
	points := make([]image.Point, 0, len(values))
	for idx, value := range values {
		x := pointX(idx, total, plot)
		y := pointY(value, minValue, maxValue, plot)
		points = append(points, image.Pt(x, y))
	}
	return points
}

func drawPolyline(img *image.RGBA, points []image.Point, c color.Color) {
	for i := 1; i < len(points); i++ {
		drawLine(img, points[i-1], points[i], c)
	}
}

func drawObservedPoints(img *image.RGBA, points []image.Point, set seriesSet) {
	normalPoint := color.RGBA{R: 18, G: 138, B: 237, A: 255}
	anomalyPoint := color.RGBA{R: 220, G: 38, B: 38, A: 255}

	for i, pt := range points {
		isAnomaly := false
		if set.expectedMax != nil && set.observed != nil && set.observed.Values[i] > set.expectedMax.Values[i] {
			isAnomaly = true
		}
		if set.expectedMin != nil && set.observed != nil && set.observed.Values[i] < set.expectedMin.Values[i] {
			isAnomaly = true
		}

		if isAnomaly {
			drawCircleBlend(img, pt, 10, color.RGBA{R: 244, G: 63, B: 94, A: 80})
			drawCircle(img, pt, 6, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			drawCircle(img, pt, 4, anomalyPoint)
			continue
		}

		drawCircle(img, pt, 4, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		drawCircle(img, pt, 3, normalPoint)
	}

}

func drawBand(img *image.RGBA, top, bottom []image.Point, fill color.RGBA) {
	if len(top) == 0 || len(bottom) == 0 || len(top) != len(bottom) {
		return
	}

	for i := 1; i < len(top); i++ {
		leftX := top[i-1].X
		rightX := top[i].X
		if rightX < leftX {
			leftX, rightX = rightX, leftX
		}

		span := rightX - leftX
		if span == 0 {
			span = 1
		}

		for x := leftX; x <= rightX; x++ {
			t := float64(x-leftX) / float64(span)
			topY := interpolate(top[i-1].Y, top[i].Y, t)
			bottomY := interpolate(bottom[i-1].Y, bottom[i].Y, t)
			if topY > bottomY {
				topY, bottomY = bottomY, topY
			}

			for y := topY; y <= bottomY; y++ {
				blendPixel(img, x, y, fill)
			}
		}
	}
}

func interpolate(a, b int, t float64) int {
	return int(float64(a) + (float64(b)-float64(a))*t)
}

func pointX(idx, total int, plot image.Rectangle) int {
	if total <= 1 {
		return plot.Min.X + plot.Dx()/2
	}
	ratio := float64(idx) / float64(total-1)
	return plot.Min.X + int(ratio*float64(plot.Dx()-1))
}

func pointY(value, minValue, maxValue float64, plot image.Rectangle) int {
	ratio := (value - minValue) / (maxValue - minValue)
	clamped := math.Max(0, math.Min(1, ratio))
	return plot.Max.Y - 1 - int(clamped*float64(plot.Dy()-1))
}

func drawLine(img *image.RGBA, from, to image.Point, c color.Color) {
	x0, y0 := from.X, from.Y
	x1, y1 := to.X, to.Y

	dx := int(math.Abs(float64(x1 - x0)))
	dy := -int(math.Abs(float64(y1 - y0)))
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy

	for {
		setPixel(img, x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func drawCircle(img *image.RGBA, center image.Point, radius int, c color.Color) {
	r2 := radius * radius
	for y := center.Y - radius; y <= center.Y+radius; y++ {
		for x := center.X - radius; x <= center.X+radius; x++ {
			dx := x - center.X
			dy := y - center.Y
			if dx*dx+dy*dy <= r2 {
				setPixel(img, x, y, c)
			}
		}
	}
}

func drawCircleBlend(img *image.RGBA, center image.Point, radius int, c color.RGBA) {
	r2 := radius * radius
	for y := center.Y - radius; y <= center.Y+radius; y++ {
		for x := center.X - radius; x <= center.X+radius; x++ {
			dx := x - center.X
			dy := y - center.Y
			if dx*dx+dy*dy <= r2 {
				blendPixel(img, x, y, c)
			}
		}
	}
}

func setPixel(img *image.RGBA, x, y int, c color.Color) {
	if image.Pt(x, y).In(img.Bounds()) {
		img.Set(x, y, c)
	}
}

func fillRect(img *image.RGBA, rect image.Rectangle, c color.RGBA) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			blendPixel(img, x, y, c)
		}
	}
}

func blendPixel(img *image.RGBA, x, y int, overlay color.RGBA) {
	if !image.Pt(x, y).In(img.Bounds()) {
		return
	}

	dst := img.RGBAAt(x, y)
	alpha := float64(overlay.A) / 255.0

	r := uint8(float64(overlay.R)*alpha + float64(dst.R)*(1-alpha))
	g := uint8(float64(overlay.G)*alpha + float64(dst.G)*(1-alpha))
	b := uint8(float64(overlay.B)*alpha + float64(dst.B)*(1-alpha))

	img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
}

func normalizeDimension(value, fallback int) int {
	if value == 0 {
		return fallback
	}
	if value < minDimension {
		return minDimension
	}
	if value > maxDimension {
		return maxDimension
	}
	return value
}

func normalizeQuality(value int) int {
	if value == 0 {
		return defaultQuality
	}
	if value < 1 {
		return 1
	}
	if value > 100 {
		return 100
	}
	return value
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func ensureJPEGFileName(fileName, fallback string) string {
	trimmed := strings.TrimSpace(fileName)
	if trimmed == "" {
		return fallback
	}

	lowered := strings.ToLower(trimmed)
	if strings.HasSuffix(lowered, ".jpeg") || strings.HasSuffix(lowered, ".jpg") {
		return trimmed
	}

	return trimmed + ".jpeg"
}

func drawText(img *image.RGBA, x, y int, label string, c color.Color) {
	d := font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(label)
}

func textWidth(label string) int {
	d := font.Drawer{Face: basicfont.Face7x13}
	return d.MeasureString(label).Ceil()
}

func shortCategoryLabel(label string) string {
	trimmed := strings.TrimSpace(label)
	if len(trimmed) <= 12 {
		return trimmed
	}
	return trimmed[:12]
}

func formatAxisValue(value float64) string {
	rounded := int64(math.Round(value))
	return formatIntWithCommas(rounded)
}

func formatIntWithCommas(n int64) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}

	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return sign + s
	}

	var b strings.Builder
	b.Grow(len(s) + len(s)/3)
	rem := len(s) % 3
	if rem == 0 {
		rem = 3
	}

	b.WriteString(s[:rem])
	for i := rem; i < len(s); i += 3 {
		b.WriteString(",")
		b.WriteString(s[i : i+3])
	}

	return sign + b.String()
}
