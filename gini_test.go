package examples

import (
	"testing"

	"github.com/irifrance/gini"
	"github.com/irifrance/gini/z"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	categories = map[string]int{"A": 1, "B": 2}
	domains    = map[string]int{"A": 3, "B": 4}
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
			So(model(g), ShouldResemble, []int{1, -2, 3, 4})
		})

		Convey(`It should not be satisfiable when "email is category A" and "email is not from domain A`, func() {
			So(isSatisfiable(g, 1, -3), ShouldBeFalse)
		})

		Convey(`It should be satisfiable when "email is category B" and "email from domain A"`, func() {
			So(isSatisfiable(g, 2, 3), ShouldBeTrue)
			So(model(g), ShouldResemble, []int{-1, 2, 3, 4})
		})

		Convey(`It should be satisfiable when "email is category B" and "email is not from domain A"`, func() {
			So(isSatisfiable(g, 2, -3), ShouldBeTrue)
			So(model(g), ShouldResemble, []int{-1, 2, -3, 4})
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
					// Untest removes assumptions since last test. In this case, -3 is removed.
					res, _ = g.Test(nil)
					So(res, ShouldNotEqual, -1)
					g.Untest()
				})
			})

			g.Untest()
		})

		Convey(`When email is from domain A, it should be classified as category A`, func() {
			So(eval(g, "A"), ShouldEqual, "A")
		})

		Convey(`When email is from domain B, it should be classified as category B`, func() {
			So(eval(g, "B"), ShouldEqual, "B")
		})

		Convey(`When email is from domain C (undefined), it should return empty (can not be classified)`, func() {
			So(eval(g, "C"), ShouldEqual, "")
		})

		Convey(`It should return a solution for "email is category A"`, func() {
			So(isSatisfiable(g, 1), ShouldBeTrue)
			So(model(g), ShouldResemble, []int{1, -2, 3, 4})
		})

		Convey(`It should return a solution for "email is category B"`, func() {
			So(isSatisfiable(g, 2), ShouldBeTrue)
			So(model(g), ShouldResemble, []int{-1, 2, 3, 4})
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

func eval(g *gini.Gini, domain string) string {
	if d, ok := domains[domain]; ok {
		for category, c := range categories {
			if isSatisfiable(g, c, d) && !isSatisfiable(g, c, -d) {
				return category
			}
		}
	}
	return ""
}

func model(g *gini.Gini) []int {
	var m []int
	for i := 1; i <= g.MaxVar().Pos().Dimacs(); i++ {
		if g.Value(z.Dimacs2Lit(i)) {
			m = append(m, i)
		} else {
			m = append(m, -i)
		}
	}
	return m
}
