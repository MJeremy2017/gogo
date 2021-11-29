package clockface

import (
	"fmt"
	"io"
	"os"
	"time"
	"github.com/quii/learn-go-with-tests/math/v6/clockface"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}

