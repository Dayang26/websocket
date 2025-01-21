package websocket

func skipSpace(s string) (rest string) {
	i := 0
	for ; i < len(s); i++ {
		if b := s[i]; b != ' ' && b != '\t' {
			break
		}
	}
	return s[i:]
}

//func nextToken(s string) (token, rest string) {
//	i := 0
//	for ; i < len(s); i++ {
//		if !isTokenObject[s[i]] {
//			break
//		}
//	}
//	return s[:i], s[i:]
//}
//
//func tokenListContainsValues(header http.Header, name, value string) bool {
//	for _, s := range header[name] {
//		for {
//			var t string
//			t, s = nextToken(skipSpace(s))
//		}
//	}
//}
