package gonf

func isBlank(r rune) bool {
	return r == ' ' || r == '\t' || isLineBreak(r)
}

func isLineBreak(r rune) bool {
	return r == '\n'
}

func isValidForKey(r rune) bool {
	// 48 == '0'
	// 57 == '9'
	isNumber := r >= 48 && r <= 57

	// 97 == 'a'
	// 122 == 'z'
	isLower := r >= 97 && r <= 122

	// 65 == 'A'
	// 90 == 'Z'
	isUpper := r >= 65 && r <= 90

	return isNumber || isLower || isUpper
}

func searchingKeyState(l *lexer) state {
	r := l.next()

	if r == '#' {
		l.stack.push(l.state)
		l.state = inCommentState
		return l.state
	}

	if r == t_EOF {
		l.finish()
		return nil
	}

	if r == '}' {
		l.emit(t_MAP_END)
		l.state = l.stack.pop()
		return l.state
	}

	if r == ']' {
		l.emit(t_ARRAY_END)
		l.state = l.stack.pop()
		return l.state
	}

	if r == '"' {
		l.ignore()
		return inQuotedKeyState
	}

	if isBlank(r) {
		l.ignore()
		return searchingKeyState
	}

	return inKeyState
}

func inQuotedKeyState(l *lexer) state {
	r := l.next()

	if r == '"' {
		l.emit(t_KEY)
		return searchingValueState
	}

	if r == '\\' {
		return inQuotedBackslashedKeyState
	}

	return inQuotedKeyState
}

func inQuotedBackslashedKeyState(l *lexer) state {
	r := l.next()

	if r == '"' || r == '\\' {
		l.backup()
		l.backup()
		l.eat()
		l.next()
	}

	return inQuotedKeyState
}

func inKeyState(l *lexer) state {
	r := l.next()

	if isBlank(r) {
		l.emit(t_KEY)
		return searchingValueState
	}

	return inKeyState
}

func searchingValueState(l *lexer) state {
	r := l.next()

	if r == '#' {
		l.stack.push(l.state)
		l.state = inCommentState
		return l.state
	}

	if isBlank(r) {
		l.ignore()
		return searchingValueState
	}

	if r == '"' {
		l.ignore()
		return inQuotedValueState
	}

	if r == '}' {
		l.emit(t_MAP_END)
		l.state = l.stack.pop()
		return l.state
	}

	if r == ']' {
		l.emit(t_ARRAY_END)
		l.state = l.stack.pop()
		return l.state
	}

	if r == '{' {
		l.emit(t_MAP_START)
		l.stack.push(l.state)
		l.state = searchingKeyState
		return l.state
	}

	if r == '[' {
		l.emit(t_ARRAY_START)
		l.stack.push(l.state)
		l.state = searchingValueState
		return l.state
	}

	return inValueState
}

func inValueState(l *lexer) state {
	r := l.next()

	if isBlank(r) {
		l.emit(t_VALUE)
		return l.state
	}

	if r == t_EOF {
		l.emit(t_VALUE)
		l.finish()
		return nil
	}

	return inValueState
}

func inQuotedValueState(l *lexer) state {
	r := l.next()

	if r == '"' {
		l.emit(t_VALUE)
		return l.state
	}

	if r == '\\' {
		return inQuotedBackslashedValueState
	}

	return inQuotedValueState
}

func inQuotedBackslashedValueState(l *lexer) state {
	r := l.next()

	if r == '"' || r == '\\' {
		l.backup()
		l.backup()
		l.eat()
		l.next()
	}

	return inQuotedValueState
}

func inCommentState(l *lexer) state {
	r := l.next()

	if isLineBreak(r) {
		l.ignore()
		l.state = l.stack.pop()
		return l.state
	}

	if r == t_EOF {
		l.finish()
		return nil
	}

	return inCommentState
}
