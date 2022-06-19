package internal

const templateComment string = `{{if .Comment}}{{printf "\n//%s" .Comment}}{{end}}`

func TrimComment(comment string) string {
	//TODO: create code that would trim a given comment
	// split the comment into words (split by space)
	// range through each of the words and count the length until it reaches the CommentWidth
	return ""
}
