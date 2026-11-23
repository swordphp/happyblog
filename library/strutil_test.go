package library

import "testing"

func TestSubStr(t *testing.T) {
	cases := []struct {
		name string
		in   string
		n    int
		want string
	}{
		{"empty", "", 5, ""},
		{"zero n", "hello", 0, ""},
		{"negative n", "hello", -1, ""},
		{"shorter than n", "hi", 5, "hi"},
		{"exact length", "hello", 5, "hello"},
		{"ascii truncate", "hello world", 5, "hello"},
		{"chinese truncate", "你好世界你好", 2, "你好"},
		{"chinese exact", "你好", 2, "你好"},
		{"mixed", "a你b好c", 3, "a你b"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := SubStr(c.in, c.n); got != c.want {
				t.Errorf("SubStr(%q,%d) = %q, want %q", c.in, c.n, got, c.want)
			}
		})
	}
}

func TestSubStrNoPanicOnShortByteString(t *testing.T) {
	// 过去的实现 content[0:200] 会在此 panic
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("SubStr panicked: %v", r)
		}
	}()
	_ = SubStr("短", 200)
}

func TestMakeBrief(t *testing.T) {
	in := "# 标题\n\n这是**正文**内容，带有 `code` 和 > 引用。"
	got := MakeBrief(in, 10)
	if r := []rune(got); len(r) > 10 {
		t.Errorf("MakeBrief returned %d runes, want <= 10", len(r))
	}
	// 不应包含被清洗掉的 markdown 噪声字符
	for _, ch := range []string{"#", "*", "`", ">", "\n"} {
		if containsStr(got, ch) {
			t.Errorf("MakeBrief result %q still contains %q", got, ch)
		}
	}
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
