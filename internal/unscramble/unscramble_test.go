package unscramble

import (
	"testing"
)

func TestIReal(t *testing.T) {
	tests := []struct {
		name        string
		scrambled   string
		unscrambled string
	}{
		{
			name:        "52 characters",
			scrambled:   "[T44A BLZC DLZE FLZG, ALZA BLZC DLZE FLZG ALZA B | ",
			unscrambled: "[T44A BLZC DLZE FLZG, ALZA BLZC DLZE FLZG ALZA B | ",
		},
		{
			name:        "Over 52 characters",
			scrambled:   "| B A BLZCZLF EZLD CZLB ZALA ,GZLF EZLD G ALZA44T[C ",
			unscrambled: "[T44A BLZC DLZE FLZG, ALZA BLZC DLZE FLZG ALZA B |C ",
		},
		{
			name:        "100 characters",
			scrambled:   "ZLB A BLZCZLF EZLD CZLB ZALA ,GZLF EZLD G ALZA44T[C- DLZE FLZG A,LZA BLZC DLZE FLZG ALZA BLZC- D |E ",
			unscrambled: "[T44A BLZC DLZE FLZG, ALZA BLZC DLZE FLZG ALZA BLZC- DLZE FLZG A,LZA BLZC DLZE FLZG ALZA BLZC- D |E ",
		},
		{
			name:        "Over 100 characters",
			scrambled:   "ZLB A BLZCZLF EZLD CZLB ZALA ,GZLF EZLD G ALZA44T[ EZLDZE FLB AZLA GZLF EZDL CZLB AZL,A GZLZC- LD -CF ",
			unscrambled: "[T44A BLZC DLZE FLZG, ALZA BLZC DLZE FLZG ALZA BLZC- DLZE FLZG A,LZA BLZC DLZE FLZG ALZA BLZC- DLZE F ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IReal(tt.scrambled)
			if result != tt.unscrambled {
				t.Errorf("IReal() =\n%s\nwant:\n%s", result, tt.unscrambled)
			}
		})
	}
}
