package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQuiz(t *testing.T) {
	quiz, err := NewQuiz(GetProblemsPath())

	Convey("When loading the quiz", t, func() {
		Convey("The quiz should be properly initialized", func() {
			So(err, ShouldBeNil)
			So(quiz.I, ShouldBeZeroValue)
			So(len(quiz.Questions), ShouldEqual, 13)
			So(quiz.Questions[0].Q, ShouldEqual, "5+5")
			So(quiz.Questions[0].A, ShouldEqual, "10")
		})

		Convey("When answering correctly", func() {
			r := quiz.CheckAnswer("10")
			So(quiz.I, ShouldEqual, 1) // Should advance the index
			So(r, ShouldBeTrue)
			So(quiz.Correct, ShouldEqual, 1)
			So(quiz.Incorrect, ShouldEqual, 0)
		})

		Convey("When answering incorrectly", func() {
			r := quiz.CheckAnswer("wrong")
			So(r, ShouldBeFalse)
			So(quiz.Correct, ShouldEqual, 1) // Did the count increase?
			So(quiz.Incorrect, ShouldEqual, 1)
		})

	})
}
