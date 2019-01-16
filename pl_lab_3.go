package main

import "net"
import "fmt"
import "bufio"
import "strings"
import "os"
import "regexp" //для регулярных выражений
import "strconv" //для преобразования типов
import "math/rand"
import "time"
import "unicode"

func get_session_key() string{
	myrand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := ""
	for i:=0;i<10;i++{
		result += string(strconv.Itoa(int(9 * myrand.Float64() + 1))[0])
	}
	return result
}

func get_hash_str() string{
	myrand := rand.New(rand.NewSource(time.Now().UnixNano()))
	// calculate initial hash string
	initial_string := ""
	for i:= 0; i<5;i++{
		initial_string += strconv.Itoa(int(6 * myrand.Float64() + 1))
	}
	return initial_string
}

func next_session_key(key string, hash string) string{
	if hash == ""{
		//handle exception
		fmt.Print("ERROR: Hash string is empty")
	}
	result := 0
    for i:=0; i<len(hash); i++ {
        temp, _:= strconv.Atoi(calc_hash(key, int([]rune(hash)[i])))
        result += temp
    }	
    result_str:= strings.Repeat("0",10) + strconv.Itoa(result)[0:10]
    return result_str[len(result_str)-10:]
}

func calc_hash(key string, val int) string{
	result := ""
	switch val{
	case 1:
		temp, _ := strconv.Atoi(key[0:5])
		temp_str := "00" + strconv.Itoa(temp % 97)
		return temp_str[len(temp_str)-2:] 
	case 2:
		for i:=1;i<len(key);i++{
			result += string(key[len(key)-i])
		}
		return result + string(key[0])
	case 3:
		return key[len(key)-5:] + key[0:5]
	case 4:
		num := 0
		for i:=1;i<9;i++{
			temp, _ := strconv.Atoi(string(key[i]))
			num += temp + 41
		}
		return strconv.Itoa(num)
	case 5:
		num := 0
		for i:=0;i<len(key);i++{
			ch := string(([]rune(key)[i]) ^ 43)
			if !unicode.IsDigit([]rune(ch)[0]) {
                ch = strconv.Itoa(int([]rune(ch)[0]))
            }
			temp, _:= strconv.Atoi(ch)
			num += temp
		}
		return strconv.Itoa(num)
	default:
		temp, _ := strconv.Atoi(key)
		return strconv.Itoa(temp + val)
	}
}

func handleRequest(conn net.Conn){
	message, _ := bufio.NewReader(conn).ReadString('\n')
	temp := strings.Split(string(message), " ")
	hash_str := temp[0]	
	previous_key := temp[1][0:len(temp[1])-1] // убираем последний символ
	fmt.Print("Initial hash: " + hash_str + " First key: " + previous_key) // лог
	next_key := next_session_key(previous_key, hash_str)
	previous_key = next_key
	fmt.Println(" Sent key: " + next_key) // лог
	conn.Write([]byte(next_key + "\n"))
	for i:=0;i<4;i++{ // цикл для 10 шагов
		message, _ := bufio.NewReader(conn).ReadString('\n')
		received_key := string(message)[0:len(string(message))-1] // убираем последний символ
		next_key = next_session_key(previous_key, hash_str)
		fmt.Print("Current key: " + next_key) // лог
		if received_key == next_key {			
			next_key = next_session_key(received_key, hash_str)
			previous_key = next_key
			conn.Write([]byte(next_key + "\n"))
			fmt.Print(" Received key: " + received_key + " Status: OK " + "Sent key: " + next_key) // лог
		}else{
			fmt.Println(" Received key: " + received_key + " ERROR" + "My cur key: " + next_key) // лог
			break
		}
	}
	conn.Close()
}

func start_server(port string){
	fmt.Println("Launching server on port: " + port)	
	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":" + port)	
	// accept connection on port
	for{
		conn, _ := ln.Accept()
		go handleRequest(conn)
	}
}

func start_client(ip_port string){
	conn, err := net.Dial("tcp", ip_port)
	if err != nil {
		// handle error
		fmt.Println("Could not connect to server")
	}else{
		hash_str := get_hash_str()
		previous_key := get_session_key()
		fmt.Println("Initial hash: " + hash_str + " First key: " + previous_key) // лог
		fmt.Fprintf(conn, hash_str + " " + previous_key + "\n")
		received_key := ""
		next_key := ""
		for i:= 0;i<4;i++{ // цикл для 10 шагов
			message, _ := bufio.NewReader(conn).ReadString('\n')
			received_key = string(message)[0:len(string(message))-1] // убираем последний символ			
			next_key = next_session_key(previous_key, hash_str)
			fmt.Print("Current key: " + next_key) // лог
			if received_key == next_key {
				next_key = next_session_key(received_key, hash_str)
				previous_key = next_key
				fmt.Fprintf(conn, next_key + "\n")
				fmt.Print(" Received key: " + received_key + " Status: OK " + "Sent key: " + next_key) // лог
			}else{
				fmt.Println(" Received key: " + received_key + " ERROR" + "My cur key: " + next_key) // лог
				break
			}
		}
		// для 10 шага (прием и сравнение без отправки
		message, _ := bufio.NewReader(conn).ReadString('\n')
		received_key = string(message)[0:len(string(message))-1] // убираем последний символ
		next_key = next_session_key(previous_key, hash_str)
		fmt.Print("Current key: " + next_key) // лог
		if received_key == next_key {
			fmt.Println(" Received key: " + received_key + " Status: OK ") // лог
		}else{
			fmt.Println(" Received key: " + received_key + " ERROR" + "My cur key: " + next_key) // лог
		}
	}
}

func main() {
	//регулярное выражение для порта
	port_regexp := regexp.MustCompile("^(([0-9]{1,4})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5]))$")
	//регулярное выражение для ip:port
	ip_port_regexp := regexp.MustCompile("^((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.){3}(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\:(([0-9]{1,4})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5]))$")
	if len(os.Args) > 1 {
		args := os.Args[1:]
		switch len(args){
		case 1:
			if port_regexp.MatchString(args[0]) {//сравнение с регулярным выражением
				start_server(args[0])
			}else{
				fmt.Println("wrong port format")
			}
		case 2:
			if ip_port_regexp.MatchString(args[0]){//сравнение с регулярным выражением	
				n, _:= strconv.Atoi(args[1])
				for i:=0;i<n;i++{ // запускаются n функций-клиентов
					//go start_client(args[0])
					start_client(args[0]) // для читабельных логов (без горутины)
				}
			}else{
				fmt.Println("wrong ip:port format")
			}
		default:
			fmt.Println("wrong number of parameters")
		}
	}else{
		fmt.Println("lack of parameters")
	}
	// для ожидания окончания работы горутины
	var stp string
    fmt.Fscan(os.Stdin, &stp)
}