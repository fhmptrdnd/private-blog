package service

import "github.com/fhmptrdnd/weather-api-test-web-based/internal/models"

// articletransform, function buat ngubah artikel
// ini contoh function sebagai tipe data
type ArticleTransform func(models.Article) models.Article

// compose, gabungin beberapa transformasi jadi satu
// contoh: compose(withtitle("judul"), withauthor("penulis"))
// implementasi rekursif: base case + recursive case
func Compose(transforms ...ArticleTransform) ArticleTransform {
	return func(a models.Article) models.Article {
		return composeHelper(a, transforms)
	}
}

// composehelper, helper function rekursif buat compose
func composeHelper(a models.Article, transforms []ArticleTransform) models.Article {
	// base case: kalo ga ada transform lagi, return artikel
	if len(transforms) == 0 {
		return a
	}
	// recursive case: apply transform pertama, lalu rekursi ke sisanya
	first := transforms[0]
	rest := transforms[1:]
	return composeHelper(first(a), rest)
}

// pipe, jalanin transformasi dari kiri ke kanan
// implementasi rekursif: base case + recursive case
func Pipe(a models.Article, transforms ...ArticleTransform) models.Article {
	// base case: kalo ga ada transform lagi, return artikel
	if len(transforms) == 0 {
		return a
	}
	// recursive case: apply transform pertama, lalu rekursi ke sisanya
	first := transforms[0]
	rest := transforms[1:]
	return Pipe(first(a), rest...)
}

// withincrementedviews, bikin function yang nambah views
// ini contoh closure: function yang return function
func WithIncrementedViews(n int) ArticleTransform {
	return func(a models.Article) models.Article {
		updated := a
		updated.Views += n
		return updated
	}
}

// withtitle, bikin function yang ganti title
func WithTitle(title string) ArticleTransform {
	return func(a models.Article) models.Article {
		updated := a
		updated.Title = title
		return updated
	}
}

// withauthor, bikin function yang ganti author
func WithAuthor(author string) ArticleTransform {
	return func(a models.Article) models.Article {
		updated := a
		updated.Author = author
		return updated
	}
}

// withcontent, bikin function yang ganti content
func WithContent(content string) ArticleTransform {
	return func(a models.Article) models.Article {
		updated := a
		updated.Content = content
		return updated
	}
}

// articleoperation, function yang bisa error
type ArticleOperation func() (models.Article, error)

// witherrorlogging, tambahin logging ke operation
func WithErrorLogging(op ArticleOperation, logMsg string) ArticleOperation {
	return func() (models.Article, error) {
		result, err := op()
		if err != nil {
			// In production: tambahkan actual logging
			_ = logMsg // placeholder
		}
		return result, err
	}
}

// retry, coba lagi kalo error
// contoh higher-order function: nerima function, return function
func Retry(op ArticleOperation, times int) ArticleOperation {
	return func() (models.Article, error) {
		var lastErr error
		for i := 0; i < times; i++ {
			result, err := op()
			if err == nil {
				return result, nil
			}
			lastErr = err
		}
	return models.Article{}, lastErr
	}
}

// map, transform setiap artikel dalam slice
// contoh: map(articles, withtitle("new title"))
// implementasi rekursif: base case + recursive case
func Map(articles []models.Article, transform ArticleTransform) []models.Article {
	// base case: kalo slice kosong, return slice kosong
	if len(articles) == 0 {
		return []models.Article{}
	}
	// recursive case: transform artikel pertama, lalu rekursi ke sisanya
	first := articles[0]
	rest := articles[1:]
	return append([]models.Article{transform(first)}, Map(rest, transform)...)
}

// filter, saring artikel berdasarkan kondisi (predicate)
// contoh: filter(articles, func(a Article) bool { return a.Views > 100 })
// implementasi rekursif: base case + recursive case
func Filter(articles []models.Article, predicate func(models.Article) bool) []models.Article {
	// base case: kalo slice kosong, return slice kosong
	if len(articles) == 0 {
		return []models.Article{}
	}
	// recursive case: cek artikel pertama, lalu rekursi ke sisanya
	first := articles[0]
	rest := articles[1:]
	filteredRest := Filter(rest, predicate)
	
	if predicate(first) {
		// kalo artikel pass kondisi, masukin ke hasil
		return append([]models.Article{first}, filteredRest...)
	}
	// kalo ga pass, skip artikel ini
	return filteredRest
}

// reduce, agregasi slice artikel jadi single value
// contoh: reduce(articles, 0, func(total int, a Article) int { return total + a.Views })
// implementasi rekursif: base case + recursive case
func Reduce[T any](articles []models.Article, initial T, reducer func(T, models.Article) T) T {
	// base case: kalo slice kosong, return nilai initial
	if len(articles) == 0 {
		return initial
	}
	// recursive case: apply reducer ke artikel pertama, lalu rekursi ke sisanya
	first := articles[0]
	rest := articles[1:]
	accumulated := reducer(initial, first)
	return Reduce(rest, accumulated, reducer)
}

// hasminviews, predicate buat filter artikel dengan minimal views
func HasMinViews(minViews int) func(models.Article) bool {
	return func(a models.Article) bool {
		return a.Views >= minViews
	}
}

// isbyauthor, predicate buat filter artikel berdasarkan author
func IsByAuthor(author string) func(models.Article) bool {
	return func(a models.Article) bool {
		return a.Author == author
	}
}

// isownedby, predicate buat filter artikel berdasarkan owner
func IsOwnedBy(ownerID string) func(models.Article) bool {
	return func(a models.Article) bool {
		return a.OwnerID == ownerID
	}
}

// parallelmap, transform setiap artikel secara concurrent
// menggunakan goroutines + channels buat parallel processing
func ParallelMap(articles []models.Article, transform ArticleTransform) []models.Article {
	if len(articles) == 0 {
		return []models.Article{}
	}

	results := make([]models.Article, len(articles))
	done := make(chan struct{}, len(articles))

	// spawn goroutine buat setiap artikel
	for i, article := range articles {
		go func(index int, a models.Article) {
			results[index] = transform(a)
			done <- struct{}{}
		}(i, article)
	}

	// tunggu semua goroutine selesai
	for i := 0; i < len(articles); i++ {
		<-done
	}
	close(done)

	return results
}

// parallelfilter, filter artikel secara concurrent
// menggunakan goroutines + channels buat parallel processing
func ParallelFilter(articles []models.Article, predicate func(models.Article) bool) []models.Article {
	if len(articles) == 0 {
		return []models.Article{}
	}

	type result struct {
		index   int
		article models.Article
		pass    bool
	}

	resultCh := make(chan result, len(articles))

	// spawn goroutine buat setiap artikel
	for i, article := range articles {
		go func(index int, a models.Article) {
			resultCh <- result{
				index:   index,
				article: a,
				pass:    predicate(a),
			}
		}(i, article)
	}

	// collect hasil yang pass predicate
	var filtered []models.Article
	for i := 0; i < len(articles); i++ {
		res := <-resultCh
		if res.pass {
			filtered = append(filtered, res.article)
		}
	}
	close(resultCh)

	return filtered
}

// pipeline, buat processing pipeline dengan channels
// articles -> transform1 -> transform2 -> ... -> hasil
func Pipeline(articles []models.Article, transforms ...ArticleTransform) <-chan models.Article {
	out := make(chan models.Article)

	go func() {
		defer close(out)
		for _, article := range articles {
			// apply semua transformasi secara berurutan
			result := article
			for _, transform := range transforms {
				result = transform(result)
			}
			out <- result
		}
	}()

	return out
}

// future, async computation pattern
// function yang return channel buat hasil async
type Future[T any] func() <-chan T

// runfuture, jalanin future dan return channel
func RunFuture[T any](computation func() T) <-chan T {
	result := make(chan T, 1)
	go func() {
		defer close(result)
		result <- computation()
	}()
	return result
}

// parallelreduce, reduce secara parallel dengan divide-and-conquer
// bagi slice jadi chunks, reduce tiap chunk parallel, lalu combine hasil
func ParallelReduce[T any](articles []models.Article, initial T, reducer func(T, models.Article) T, combiner func(T, T) T) T {
	if len(articles) == 0 {
		return initial
	}

	// kalo cuma sedikit, pake reduce biasa
	if len(articles) < 10 {
		return Reduce(articles, initial, reducer)
	}

	// bagi jadi 2 bagian
	mid := len(articles) / 2
	left := articles[:mid]
	right := articles[mid:]

	// process parallel
	leftCh := make(chan T, 1)
	rightCh := make(chan T, 1)

	go func() {
		leftCh <- ParallelReduce(left, initial, reducer, combiner)
		close(leftCh)
	}()

	go func() {
		rightCh <- ParallelReduce(right, initial, reducer, combiner)
		close(rightCh)
	}()

	// combine hasil
	leftResult := <-leftCh
	rightResult := <-rightCh

	return combiner(leftResult, rightResult)
}

// fanout, distribute artikel ke multiple workers
// setiap worker jalanin transform function
func FanOut(articles []models.Article, numWorkers int, transform ArticleTransform) []models.Article {
	if len(articles) == 0 || numWorkers <= 0 {
		return []models.Article{}
	}

	jobs := make(chan models.Article, len(articles))
	results := make(chan models.Article, len(articles))

	// spawn workers
	for w := 0; w < numWorkers; w++ {
		go func() {
			for article := range jobs {
				results <- transform(article)
			}
		}()
	}

	// kirim jobs
	go func() {
		for _, article := range articles {
			jobs <- article
		}
		close(jobs)
	}()

	// collect hasil
	output := make([]models.Article, 0, len(articles))
	for i := 0; i < len(articles); i++ {
		output = append(output, <-results)
	}
	close(results)

	return output
}
