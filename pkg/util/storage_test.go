package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
)

func Test_filterFilesToDelete(t *testing.T) {
	type args struct {
		maxDays   int
		maxMonths int
		files     []File
	}
	tests := []struct {
		name string
		args args
		want []File
	}{
		{
			name: "Test_filterFilesToDelete",
			args: args{maxDays: 3, maxMonths: 6, files: []File{
				{Name: "2021.12.01/TEZOS_HANGZHOUNET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"}, // Delete
				{Name: "2021.12.02/TEZOS_HANGZHOUNET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"},
				{Name: "2021.12.02/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"},    // Delete
				{Name: "2021.12.02/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.rolling"}, // Delete
				{Name: "2022.01.02/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"},
				{Name: "2022.02.02/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"},
				{Name: "2022.03.02/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"},
				{Name: "2022.04.02/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"},
				{Name: "2022.05.02/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.rolling"},
				{Name: "2022.05.21/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"}, 
				{Name: "2022.05.21/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.rolling"},//Delete
				{Name: "2022.05.22/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.rolling"}, //Delete
				{Name: "2022.05.26/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"},
				{Name: "2022.05.26/TEZOS_ITHACANET_2022-01-25T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"}, //Delete
				{Name: "2022.05.26/TEZOS_ITHACANET_2022-01-25T15:00:00Z-BLHBNkniS9kE6SME1jQRSp3xLCiAiqsyUTXXS7LdtZjcF4ZXVEX-589921.rolling"},
				{Name: "2022.05.26/TEZOS_JAKARTANET_2022-04-27T15:00:00Z-BLBVjuvySZoQ6AW5BzGbnrSsYqjd6FBZiBdHNKicmUgdVsSD6Wu-153024.full"},
				{Name: "2022.05.26/TEZOS_JAKARTANET_2022-04-27T15:00:00Z-BLCaDAV2CJgEqR18YDHUCjZaCxWV7xJmtbT8d2GNMzpghVednkb-153022.rolling"},
				{Name: "2022.05.26/TEZOS_MAINNET-BL4p3YRfxhiQP16PsuzFBbph8QqVcNN3qu42r5JgNgdaw3xW81g-2396664.rolling"},
				{Name: "2022.05.26/TEZOS_MAINNET-BLFkccePHzNuCPEBSQPuxiwB95R45RdvjEebJDC8GTQJfSwfy1E-2396725.full"},
			}},
			want: []File{
				{Name: "2021.12.01/TEZOS_HANGZHOUNET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"},
				{Name: "2021.12.02/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"},
				{Name: "2021.12.02/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.rolling"},
				{Name: "2022.05.21/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.rolling"},
				{Name: "2022.05.22/TEZOS_ITHACANET_2022-01-22T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.rolling"},
				{Name: "2022.05.26/TEZOS_ITHACANET_2022-01-25T15:00:00Z-BKpkmdGCx8D9KAUAYJrrrmFqgamcwZWFYo2W4KiyEP4PCBJQrsC-589926.full"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterFilesToDelete(tt.args.maxDays, tt.args.maxMonths,
				tt.args.files, time.Date(2022, 6, 26, 0, 0, 0, 0, time.UTC))
			fmt.Printf("%d\n", len(got))

			if len(got) != 6 {
				t.Errorf("filterFilesToDelete() = %v, want %v", len(got), 6)
			}

			lo.ForEach(tt.want, func(f File, i int) {
				if !lo.Contains(got, f) {
					t.Errorf("got: %v, want contains: %v", got, f)
				}
			})
		})
	}
}
