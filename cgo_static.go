//go:build !duckdb_use_lib && (darwin || (linux && (amd64 || arm64)) || (freebsd && amd64) || (windows && amd64))

package duckdb

/*
#cgo LDFLAGS:  
#cgo darwin,amd64 LDFLAGS: -lduckdb -lc++ -L${SRCDIR}/deps/darwin_amd64
#cgo darwin,arm64 LDFLAGS: -lduckdb -lc++ -L${SRCDIR}/deps/darwin_arm64
#cgo linux,amd64 LDFLAGS: -lduckdb -lstdc++ -lm -ldl -L${SRCDIR}/deps/linux_amd64
#cgo linux,arm64 LDFLAGS: -lduckdb -lstdc++ -lm -ldl -L${SRCDIR}/deps/linux_arm64
#cgo windows,amd64 LDFLAGS: -lduckdb_99 -lduckdb_98 -lduckdb_97 -lduckdb_96 -lduckdb_95 -lduckdb_94 -lduckdb_93 -lduckdb_92 -lduckdb_91 -lduckdb_90 -lduckdb_9 -lduckdb_89 -lduckdb_88 -lduckdb_87 -lduckdb_86 -lduckdb_85 -lduckdb_84 -lduckdb_83 -lduckdb_82 -lduckdb_81 -lduckdb_80 -lduckdb_8 -lduckdb_79 -lduckdb_78 -lduckdb_77 -lduckdb_76 -lduckdb_75 -lduckdb_74 -lduckdb_73 -lduckdb_72 -lduckdb_71 -lduckdb_70 -lduckdb_7 -lduckdb_69 -lduckdb_68 -lduckdb_67 -lduckdb_66 -lduckdb_65 -lduckdb_64 -lduckdb_63 -lduckdb_62 -lduckdb_61 -lduckdb_60 -lduckdb_6 -lduckdb_59 -lduckdb_58 -lduckdb_57 -lduckdb_56 -lduckdb_55 -lduckdb_54 -lduckdb_53 -lduckdb_52 -lduckdb_51 -lduckdb_50 -lduckdb_5 -lduckdb_49 -lduckdb_48 -lduckdb_47 -lduckdb_46 -lduckdb_45 -lduckdb_44 -lduckdb_43 -lduckdb_42 -lduckdb_41 -lduckdb_40 -lduckdb_4 -lduckdb_39 -lduckdb_38 -lduckdb_37 -lduckdb_36 -lduckdb_35 -lduckdb_34 -lduckdb_33 -lduckdb_32 -lduckdb_31 -lduckdb_30 -lduckdb_3 -lduckdb_29 -lduckdb_28 -lduckdb_27 -lduckdb_26 -lduckdb_25 -lduckdb_24 -lduckdb_23 -lduckdb_22 -lduckdb_21 -lduckdb_20 -lduckdb_2 -lduckdb_19 -lduckdb_18 -lduckdb_17 -lduckdb_16 -lduckdb_155 -lduckdb_154 -lduckdb_153 -lduckdb_152 -lduckdb_151 -lduckdb_150 -lduckdb_15 -lduckdb_149 -lduckdb_148 -lduckdb_147 -lduckdb_146 -lduckdb_145 -lduckdb_144 -lduckdb_143 -lduckdb_142 -lduckdb_141 -lduckdb_140 -lduckdb_14 -lduckdb_139 -lduckdb_138 -lduckdb_137 -lduckdb_136 -lduckdb_135 -lduckdb_134 -lduckdb_133 -lduckdb_132 -lduckdb_131 -lduckdb_130 -lduckdb_13 -lduckdb_129 -lduckdb_128 -lduckdb_127 -lduckdb_126 -lduckdb_125 -lduckdb_124 -lduckdb_123 -lduckdb_122 -lduckdb_121 -lduckdb_120 -lduckdb_12 -lduckdb_119 -lduckdb_118 -lduckdb_117 -lduckdb_116 -lduckdb_115 -lduckdb_114 -lduckdb_113 -lduckdb_112 -lduckdb_111 -lduckdb_110 -lduckdb_11 -lduckdb_109 -lduckdb_108 -lduckdb_107 -lduckdb_106 -lduckdb_105 -lduckdb_104 -lduckdb_103 -lduckdb_102 -lduckdb_101 -lduckdb_100 -lduckdb_10 -lduckdb_1 -lduckdb_0 -lstdc++ -lm -L${SRCDIR}/deps/windows_amd64
#cgo freebsd,amd64 LDFLAGS: -lduckdb -lstdc++ -lm -ldl -L${SRCDIR}/deps/freebsd_amd64
#include <duckdb.h>
*/
import "C"
