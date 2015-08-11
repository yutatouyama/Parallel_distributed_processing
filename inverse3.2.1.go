package main

import (
    fmt "fmt"
    time "time"
    sync "sync"
    "runtime"
    //"bufio"
    rand "math/rand"
)

const SIZE int = 1000

type Matrix struct {
	data [SIZE][SIZE]float64
}

type Bvector struct {
	data [SIZE]float64
} 

var epsilon float64 = 1e-10

func main() {

	var err int 
	var I, A Matrix

	err = 0

	fmt.Println(SIZE,"次正方行列の逆行列を求めるプログラムです\n")

	runtime.GOMAXPROCS(runtime.NumCPU())

	//A = MatrixInput();
	A = MatrixCreate();
	fmt.Println("\n生成された行列は\n");
	MatrixPrint(A);

	s := time.Now()
	I = MatrixInverse(A, &err);
	e := time.Now().Sub(s)
    fmt.Println("逆行列生成時間",e)
	if false {
		fmt.Println("\n逆行列は存在しません\n")
		return
	} else {
		fmt.Println("\n逆行列は\n");
		MatrixPrint(I);
	}
}

func MatrixCreate() Matrix {
	var X Matrix

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			X.data[i][j] = float64(rand.Intn(10)) - float64(rand.Intn(50))
		}
	}

	return X
}


/*
func MatrixInput() Matrix {
	var X Matrix 
	var i, j int

	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
      		fmt.Println("%d 次正方行列Aの (%d,%d) 成分を入力してください\n", SIZE, i+1, j+1);
      		scanner := bufio.NewScanner("%lf", &X.data[i][j]);
    	}
	}
	return X
}
*/

func MatrixInverse(X Matrix, error *int) Matrix {
	var R, LUM Matrix
	var toggle, err int
	var perm, b, x Bvector
	var wg sync.WaitGroup

	err = 0
	R = MatrixDuplicate(X)
	
	a := make(chan *Matrix, runtime.NumCPU())
	wg.Add(1)
	go func () {		
		defer wg.Done()
		MatrixDecompose(X, &perm, &toggle, &err)
		a <- &R 
	} ()
	LUM = *(<- a)
	wg.Wait()

	if false {
		*error = 1
		R = MatrixCreate()
	} else {
		for i := 0; i < SIZE; i++ {
      		for j := 0; j < SIZE; j++ {
				if float64(i) == perm.data[j] {
					b.data[j] = 1
				} else {
					b.data[j] = 0
				}
      		}

      		x = HelperSolve(LUM, b)
      		
      		for j := 0; j < SIZE; j++ {
      			R.data[j][i] = x.data[j];
      		}   		 
    	}
	}	
	return R
}

func  MatrixDuplicate(X Matrix) Matrix {
	var Y Matrix

	for i := 0; i < SIZE; i++ {
    	for j := 0; j < SIZE; j++ {
      		Y.data[i][j] = X.data[i][j];
    	}
  	}
  	return Y
}

 func MatrixDecompose(X Matrix, P *Bvector, toggle *int, error *int) Matrix {
	var R Matrix
	R = MatrixDuplicate(X)
	var row int
	var colMax, rowPtr, temp float64

	for i := 0; i < SIZE; i++ {
    	(*P).data[i] = float64(i)
  	}
  	*toggle = 1

  	for j := 0; j < SIZE-1; j++ {
    	colMax = abs(R.data[j][j])
    	row = j
    	for i := j+1; i < SIZE; i++ {
      		if R.data[i][j] > colMax {
				colMax = R.data[i][j]
				row = i
      		}
    	}
    	if row != j {
      		for k := 0; k < SIZE; k++ {
				rowPtr = R.data[row][k]
				R.data[row][k] = R.data[j][k]
				R.data[j][k] = rowPtr
      		}
      	temp = (*P).data[row]
      	(*P).data[row] = (*P).data[j]
      	(*P).data[j] = temp
      	*toggle = -(*toggle)
    	}
    	if abs(R.data[j][j]) < epsilon {
    	  *error = 1
    	} else {
      		for i := j+1; i < SIZE; i++ {
				R.data[i][j] /= R.data[j][j]
				for k := j+1; k < SIZE; k++ {
	  				R.data[i][k] -= R.data[i][j] * R.data[j][k]
				}
      		}
    	}
  	}
  return R

}

func HelperSolve(LU Matrix, b Bvector) Bvector {

	var sum float64
	var	x Bvector

	for i := 0; i < SIZE; i++ {
    	x.data[i] = b.data[i]
  	}
	for i := 1; i < SIZE; i++ {
    	sum = x.data[i]
    	for j := 0; j < i; j++ {
      		sum -= LU.data[i][j] * x.data[j]
    	}
    x.data[i] = sum
  	}
  	x.data[SIZE-1] /= LU.data[SIZE-1][SIZE-1]
  		for i := SIZE-2; i >= 0; i-- {
    		sum = x.data[i]
    		for j := i+1; j < SIZE; j++ {
      			sum -= LU.data[i][j] * x.data[j]
    		}
   		 x.data[i] = sum / LU.data[i][i]
  		}
 
  return x
}

func MatrixPrint(X Matrix) int {

	for i := 0; i < SIZE; i++ {
    	for j := 0; j < SIZE; j++ {
      		fmt.Print(X.data[i][j],", ")
    	}	
    fmt.Println("\n")
  }
 
  return 0

}

func abs(x float64) float64 {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0 // return correctly abs(-0)
	}
	return x
}


