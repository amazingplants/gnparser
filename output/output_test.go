package output

import (
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	grm "gitlab.com/gogna/gnparser/grammar"
)

var _ = Describe("Output", func() {

})

var _ = Describe("Private Functions", func() {
	Describe("prepareWarnings", func() {
		It("sorts incoming data consistently", func() {
			ws := []grm.Warning{
				grm.YearParensWarn,
				grm.YearCharWarn,
				grm.CharBadWarn,
				grm.YearQuestionWarn,
				grm.TailWarn,
			}
			res := prepareWarnings(ws)
			output := []Warning{
				{
					Quality: 3,
					Message: "Unparsed tail",
				},
				{
					Quality: 2,
					Message: "Non-standard characters in canonical",
				},
				{
					Quality: 2,
					Message: "Year with latin character",
				},
				{
					Quality: 2,
					Message: "Year with parentheses",
				},
				{
					Quality: 2,
					Message: "Year with question mark",
				},
			}
			Expect(res).To(Equal(output))
		})
	})
})

func randIntSlice(sl []int) []int {
	res := make([]int, len(sl))
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i, v := range r.Perm(len(sl)) {
		res[i] = sl[v]
	}
	return res
}
