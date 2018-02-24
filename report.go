package peepingJim

import (
	"fmt"
	"log"
	"os"
)

//BuildReport takes a hashmap and builds an html file that will be written
//to the file system
func BuildReport(db []map[string]string, outFile string) {
	var liveMarkup string
	for _, d := range db {
		liveMarkup += fmt.Sprintf(`<tr><td class='img'><a href='%s' target='_blank'><img src='%s' onerror="this.parentNode.parentNode.innerHTML='No image available.';" /></a></td><td class='head'><a href='%s' target='_blank'>%s</a> (<a href='%s' target='_blank'>source</a>)<br /><pre>%s</pre></td></tr>`, d["imgPath"], d["imgPath"], d["url"], d["url"], d["srcPath"], d["headers"])
	}
	htmlBody := fmt.Sprintf(`<!doctype html>
<head>
<style>
table, td, th {border: 1px solid black;border-collapse: collapse;padding: 5px;font-size: .9em;font-family: tahoma;}
table {width: 100%%;table-layout: fixed;min-width: 1000px;}
td.img {width: 40%%;}
img {width: 100%%;}
td.head {vertical-align: top;word-wrap: break-word;}
</style>
</head>
<body>
<table>
%s
</table>
</body>
</html>`, liveMarkup)
	file, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(htmlBody)
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}
