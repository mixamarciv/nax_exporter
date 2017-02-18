package app_fnc

//запускает cnt_workers gorutin и выполняет в них fnc для обработки строк из arrstrjobs
func Exec_go_workers(fnc func(string, interface{}) string, arrstrjobs []string, params interface{}, cnt_workers int) []string {
	cnt_w := cnt_workers //максимальное количество горутин

	cnt_j := len(arrstrjobs) //количество заданий
	if cnt_j < cnt_w {       //ограничеваем количество горутин
		cnt_w = cnt_j //должно быть не больше чем всего заданий
	}

	jobs := make(chan string, cnt_w)
	results := make(chan string, cnt_w)

	//Printf("\nExec_go_workers:  cnt_w %d  cnt_j %d", cnt_w, cnt_j)
	//стартуем горутины
	for w := 1; w <= cnt_w; w++ {
		go run_go_worker(fnc, w, jobs, params, results)
	}

	//отправляем задания
	for j := 0; j < cnt_j; j++ {
		jobs <- arrstrjobs[j]
	}
	close(jobs)

	data := make([]string, cnt_j)
	for j := 0; j < cnt_j; j++ {
		res := <-results
		data[j] = res
	}

	return data
}

func run_go_worker(fnc func(string, interface{}) string, id int, jobs <-chan string, params interface{}, results chan<- string) {
	for j := range jobs {
		s := Sprintf("#worker%d %s\n", id, j)
		s += fnc(j, params)
		results <- s
	}
}
