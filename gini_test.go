package examples

import (
	"testing"

	"github.com/irifrance/gini"
	"github.com/irifrance/gini/z"
	. "github.com/smartystreets/goconvey/convey"
)

//
// Giving the email filtering criteria:
//
// Email is from domain A, it should be classified as category A
// Email is from domain B, it should be classified as category B
//
// Literals:
//   1: email is category A
//   2: email is category B
//   3: email is from domain A
//   4: email is from domain B
//
// Clauses:
//   1 -> 3 = !1 | 3 = [-1, 3]
//   2 -> 4 = !2 | 4 = [-2, 4]
//
func TestGini(t *testing.T) {
	Convey("Giving a sat solver with the clauses", t, func() {
		g := gini.New()

		addClause(g, -1, 3)
		addClause(g, -2, 4)

		Convey(`It should be satisfiable when "email is category A" and "email from domain A"`, func() {
			So(isSatisfiable(g, 1, 3), ShouldBeTrue)
		})

		Convey(`It should not be satisfiable when "email is category A" and "email is not from domain A`, func() {
			So(isSatisfiable(g, 1, -3), ShouldBeFalse)
		})

		Convey(`It should be satisfiable when "email is category B" and "email from domain A"`, func() {
			So(isSatisfiable(g, 2, 3), ShouldBeTrue)
		})

		Convey(`It should be satisfiable when "email is category B" and "email is not from domain A"`, func() {
			So(isSatisfiable(g, 2, -3), ShouldBeTrue)
		})

		Convey(`It should be satisfiable when "email is category A"`, func() {
			assume(g, 1)

			res, _ := g.Test(nil)
			So(res, ShouldNotEqual, -1)

			Convey(`When "email from domain A" is also true`, func() {
				assume(g, 3)
				Convey(`It should be still satisfiable`, func() {
					res, _ := g.Test(nil)
					So(res, ShouldNotEqual, -1)
					g.Untest()
				})
			})

			Convey(`When "email from domain A" is not true`, func() {
				assume(g, -3)
				Convey(`It should not be satisfiable`, func() {
					res, _ := g.Test(nil)
					So(res, ShouldEqual, -1)
					g.Untest()
				})
			})

			g.Untest()
		})
	})
}

func addClause(g *gini.Gini, lits ...int) {
	for _, lit := range lits {
		g.Add(z.Dimacs2Lit(lit))
	}
	g.Add(0)
}

func assume(g *gini.Gini, assumptions ...int) {
	for _, assumption := range assumptions {
		g.Assume(z.Dimacs2Lit(assumption))
	}
}

func isSatisfiable(g *gini.Gini, assumptions ...int) bool {
	assume(g, assumptions...)
	return g.Solve() == 1
}
