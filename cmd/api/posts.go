package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/sikozonpc/social/internal/store"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--createPostHandler")
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: Change after auth
		UserID: 202,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("---get post handler")
	// post := getPostFromCtx(r)
	// fmt.Println("---get post handler", post)
	// comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	// if err != nil {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }

	// post.Comments = comments

	// if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }
	fmt.Println("---get post handler")
	idParam := chi.URLParam(r, "postID")
	fmt.Println("---a")
	id, err := strconv.ParseInt(idParam, 10, 64)
	fmt.Println("---b")
	if err != nil {
		fmt.Println("---c")
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		fmt.Println("---d")
		return
	}
	fmt.Println("---before reading context")
	ctx := r.Context()
	fmt.Println("---after reading context")
	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("---END get post handler")
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--delete post handler")
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Posts.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	fmt.Println("--postsContextMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		fmt.Println("--adAD")
		id, err := strconv.ParseInt(idParam, 10, 64)
		fmt.Println("--BABAB")
		if err != nil {
			fmt.Println("--internal server error")
			app.internalServerError(w, r, err)
			return
		}
		fmt.Println("--after error")
		ctx := r.Context()

		post, err := app.store.Posts.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	fmt.Println("--get post")
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
