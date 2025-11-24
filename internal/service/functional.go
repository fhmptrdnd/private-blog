package service

import "github.com/fhmptrdnd/weather-api-test-web-based/internal/models"

// articletransform, function buat ngubah artikel
// ini contoh function sebagai tipe data
type ArticleTransform func(models.Article) models.Article

// compose, gabungin beberapa transformasi jadi satu
// contoh: compose(withtitle("judul"), withauthor("penulis"))
func Compose(transforms ...ArticleTransform) ArticleTransform {
	return func(a models.Article) models.Article {
		result := a
		for _, t := range transforms {
			result = t(result)
		}
		return result
	}
}

// pipe, jalanin transformasi dari kiri ke kanan
func Pipe(a models.Article, transforms ...ArticleTransform) models.Article {
	result := a
	for _, t := range transforms {
		result = t(result)
	}
	return result
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
