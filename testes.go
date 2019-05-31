package main

import(

	"fmt"
	"encoding/json"
	"os"
	
)



type response1 struct{
	Page int
	Fruits []string
}

func main(){

	
	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))
	intB, _ := json.Marshal(1)
	fmt.Println(string(intB))
	
	fltB, _ := json.Marshal(2.34)
	fmt.Println(string(fltB))
	
	strB, _ := json.Marshal("gopher")
	fmt.Println(string(strB))
	
	slcD := []string{"apple","peach","pear"}
	slcM, _ := json.Marshal(slcD)
	fmt.Println(string(slcM))
	
	mapD := map[string]int{"apple":5,"orange":12}
	mapM, _ := json.Marshal(mapD)
	fmt.Println(string(mapM))
	

	res1D := &response1{
		Page: 1,
		Fruits: []string{"apple","peasch","pear","orange"},
	}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))

	byt := []byte( `{"num":6.13,"strs":["a","b","c"],"num2":55}` )

	var dat map[string]interface{}

	if err := json.Unmarshal(byt, &dat)
	err != nil{
		panic(err)
	}
	fmt.Println(dat)

	num := dat["num"].(float64)
	fmt.Println(num)

	num2 := dat["num2"].(float64)
	fmt.Println(num2)

	strs := dat["strs"].([]interface{})
	str1:= strs[0].(string) 
	fmt.Println(str1)

	enc := json.NewEncoder(os.Stdout)
	d := map[string]int{"apple":5, "lettuce":8 } 
	enc.Encode(d)
}