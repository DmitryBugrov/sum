package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"math"
)



type number_string struct {
	name  string
	rubl string
	kop string
}

var number =map[int64]number_string{
	1: {name: "один", rubl: "рубль", kop: "копейка"},
	2: {name: "два", rubl: "рубля", kop: "копейки"},
	3: {name: "три", rubl: "рубля", kop: "копейки"},
	4: {name: "четыре", rubl: "рубля", kop: "копейки"},
	5: {name: "пять", rubl: "рублей", kop: "копеек"},
	6: {name: "шесть", rubl: "рублей", kop: "копеек"},
	7: {name: "семь", rubl: "рублей", kop: "копеек"},
	8: {name: "восемь", rubl: "рублей", kop: "копеек"},
	9: {name: "девять", rubl: "рублей", kop: "копеек"},
	0: {name: "", rubl: "рублей", kop: "копеек"},
	10: {name: "десять", rubl: "рублей", kop: "копеек"},
	11: {name: "одиндцать", rubl: "рублей", kop: "копеек"},
	12: {name: "двенадцать", rubl: "рублей", kop: "копеек"},
	13: {name: "тринадцать", rubl: "рублей", kop: "копеек"},
	14: {name: "четырнадцать", rubl: "рублей", kop: "копеек"},
	15: {name: "пятнадцать", rubl: "рублей", kop: "копеек"},
	16: {name: "шестнадцать", rubl: "рублей", kop: "копеек"},
	17: {name: "семнадцать", rubl: "рублей", kop: "копеек"},
	18: {name: "восемнадцать", rubl: "рублей", kop: "копеек"},
	19: {name: "девятнадцать", rubl: "рублей", kop: "копеек"},
	
}
var ten = map[int64]string{
	1: "десять",
	2: "двадцать",
	3: "тридцать",
	4: "сорок",
	5: "пятьдесят",
	6: "шестьдесят",
	7: "семьдесят",
	8: "восемьдесят",
	9: "девяносто",
}

var hundred = map[int64]string{
	1: "сто",
	2: "двести",
	3: "триста",
	4: "четыреста",
	5: "пятьсот",
	6: "шестьсот",
	7: "семьсот",
	8: "восемьсот",
	9: "девятьсот",
}

func Parsing (w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	input_string := params["rubl"]
	//заменяем во входной строке разделитель "," на "."
	fSum, err :=strconv.ParseFloat(strings.Replace(input_string,",",".",-1),64)
	//отделяем рубли
	int, frac := math.Modf(fSum)
	iRubl :=int64(int)
	//от копеек
	fKop:=frac*100
	_,frac=math.Modf(fKop)
	if frac >= 0.5 {
		fKop = math.Ceil(fKop)
	} else {
		frac = math.Floor(fKop)
	}
	iKop:=int64(fKop)
	
	//выходная строка
	sum := ""
	pos:=1
	//ограничение
	if iRubl>999999999 {err=errors.New("Ошибка, число слишком большое!!!\n")}
	if iKop>99 {err=errors.New("Ошибка в передаваемой строке!!!\n")}
	if err == nil {
			for iRubl>0 {
				digit:=iRubl%10
				switch pos {
					//обрабатываем каждый первый разряд в тройке
				case 1,4,7:
					if iRubl%100/10==1 {
						digit=iRubl%100
						}
						
					//добавляем "тысячи", если необходимо
					if pos==4 && iRubl%1000>0 {
						switch digit  {
						case 1:
							sum=" тысяча"+sum
						case 2,3,4:
							sum=" тысячи"+sum
						default:
							sum=" тысячь"+sum
						}
					}
					//добавляем "миллионы", если необходимо
					if pos==7 {
						switch digit  {
						case 1:
							sum=" миллион"+sum
						case 2,3,4:
							sum=" миллиона"+sum
						default:
							sum=" миллионов"+sum
						}
					}
					
					// добавляем "рублей" в требуемой форме
					if sum=="" {
						sum=" "+number[digit].rubl
						}
					
					//корректировка для тысяч
					if pos==4 && digit==1 {
						sum=" одна"+sum
					} else {
						if pos==4 && digit==2 {
							sum=" две"+sum
						} else {
						if digit!=0 {
							sum=" "+number[digit].name+sum
							}
						}
					}
					
					pos++
					if iRubl%100/10==1 { 
						pos++
						iRubl=iRubl/10
					}
				case 2,5,8:
					if digit!=0 {
						sum=" "+ten[digit]+sum;
					}
					pos++
				case 3,6,9:
					if digit!=0 {
						sum=" "+hundred[digit]+sum;
					}
					pos++
				}
				iRubl=iRubl/10
				
			}
			switch {
				case iKop>=3 && iKop<=19:
					sum=sum+ " "+number[iKop].name+" "+number[iKop].kop	
				case iKop%10==0:
					sum=sum+ " "+ten[iKop/10]+" "+number[0].kop
				case iKop%10==1:
					sum=sum+ten[iKop/10]+" одна "+number[1].kop	
				case iKop%10==2:
					sum=sum+ten[iKop/10]+" две "+number[2].kop	
				default:
					sum=sum+" "+ten[iKop/10]+" "+number[iKop%10].name+" "+number[iKop%10].kop

			}

	}
	if err == nil {
			fmt.Fprintf(w, "%s\n", sum)
		} else {
			fmt.Fprintf(w, "Ошибка в передаваемой строке: %s\n%s \n\n\n Используйте /sum/Rubl,Kop \n \t где: Rubl - сумма в рублях, до 999 млн. руб., Kop - сумма копеек от 0 до 99", input_string,err)
		}
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
         fmt.Fprint(w, "Используйте /sum/Rubl,Kop \n \t где: Rubl - сумма в рублях, до 999 млн. руб., Kop - сумма копеек от 0 до 99")
    }

func main() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/sum/{rubl:[0-9,.]+}", Parsing) 
	rtr.NotFoundHandler = http.HandlerFunc(errorHandler)
	err := http.ListenAndServe(":1234", rtr)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
