package mojibake

import (
	"testing"
)

var data = []string{
	"(Σ╕ÇΦê¼σ░ÅΦ¬¼) [µí£σ║¡Σ╕Çµ¿╣] GOSICK 1∩╜₧6σ╖╗∩╝ïGOSICKs 1∩╜₧3σ╖╗ (Θ¥Æτ⌐║µûçσ║½txtσ╜óσ╝Å µî┐τ╡╡Σ╗ÿπüì)(090518µáíµ¡ú)",
	"GOSICK τ¼¼1σ╖╗.txt",
	"πâíπâó.txt",
	"(êΩö╩Å¼Éα) [ôⁿè╘Élè╘] ôdögÅùé╞É┬ÅtÆj æµ01è¬ (É┬ï≤ò╢î╔æ╬ë₧txt ò/ÄåüEæ}èGòt)(ìZÉ│10-02-10).txt",
	"[ò╩ì√Å¡öNâ}âKâWâô9îÄìå][µ_ÄRæn]Éiîéé╠ïÉÉlü@æµ47üE48ÿb",
	"\x93d\x94g\x8f\x97\x82\xc6\x90\xc2\x8ft\x92j\x81@\x91\xe6\x82P\x8a\xaa",
	"\x93\xfc\x8a\xd4\x90l\x8a\xd4",
}

var decoded = []string{
	"(一般小説) [桜庭一樹] GOSICK 1～6巻＋GOSICKs 1～3巻 (青空文庫txt形式 挿絵付き)(090518校正)",
	"GOSICK 第1巻.txt",
	"メモ.txt",
	"(一般小説) [入間人間] 電波女と青春男 第01巻 (青空文庫対応txt \x00紙・挿絵付)(校正10-02-10).txt",
	"[別冊少年マガジン9月号][訐山創]進撃の巨人　第47・48話",
	"電波女と青春男　第１巻",
	"入間人間",
}

func TestDecode(t *testing.T) {
	test := make([]string, len(data))

	test[0] = MustDecode(data[0], CP473)
	test[1] = MustDecode(data[1], CP473)
	test[2] = MustDecode(data[2], CP473)
	test[3] = MustDecode(data[3], CP473, CP932)
	test[4] = MustDecode(data[4], CP473, CP932)
	test[5] = MustDecode(data[5], CP932)
	test[6] = MustDecode(data[6], CP932)

	for i, s := range test {
		if s != decoded[i] {
			t.Errorf("'%s' should be '%s'\n", s, decoded[i])
		}
	}
}
