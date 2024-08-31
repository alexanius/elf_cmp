# elf_cmp

This is a tool for comparing two elf files. It can show you the similar and
different symbols and make a diff of functions.

## HTML report

The recomended report type is html. It can be simply generated with:

```
$ elf_cmp -html file1 file2

$ firefox report/index.html
```

## Default launch

The default working mode is launching util with two binaries:

```
$ elf_cmp file1 file2
+------------------------------------+--------------------------------------------------------------------+
| A                                  | testdata/main.old                                                  |
| B                                  | ../edu/beta/bin/beta                                               |
+------------------------------------+----------------------+----------------------+----------------------+
|                                    | A                    | B                    | Diff                 |
+------------------------------------+----------------------+----------------------+----------------------+
| General info                                                                                            |
+------------------------------------+----------------------+----------------------+----------------------+
| Type                               | Executable file      | Shared object file   | !                    |
| Debug info                         | yes                                         |                      |
| Sections                           | 23                   | 40                   | !                    |
| Size                               | 1921273              | 791856               | -58.78%              |
+------------------------------------+----------------------+----------------------+----------------------+
| Sections                                                                         |                      |
+---------------+--------------------+----------------------+----------------------+----------------------+
| Instr         | .text              | 536809               | 211031               | -60.69%              |
|               | .plt               |                      | 1504                 |                      |
|               | .plt.got           |                      | 8                    |                      |
|               | .fini              |                      | 13                   |                      |
|               | .init              |                      | 27                   |                      |
|               | Total              | 536809               | 212583               | -60.40%              |
+---------------+--------------------+----------------------+----------------------+----------------------+
| User data     | .data              | 16272                | 13728                | -15.63%              |
|               | .bss               | 188056               | 280                  | -99.85%              |
|               | .rodata            | 246849               | 22778                | -90.77%              |
|               | Total              | 451177               | 36786                | -91.85%              |
+---------------+--------------------+----------------------+----------------------+----------------------+
| Go data       | .gopclntab         | 415064               |                      |                      |
|               | .noptrdata         | 21728                |                      |                      |
|               | .noptrbss          | 15920                |                      |                      |
|               | .typelink          | 1452                 |                      |                      |
|               | .itablink          | 176                  |                      |                      |
|               | .gosymtab          | 0                    |                      |                      |
|               | Total              | 454340               | 0                    | -100.00%             |
+---------------+--------------------+----------------------+----------------------+----------------------+
| Compiler data | .go.buildinfo      | 304                  |                      |                      |
|               | .note.go.buildid   | 100                  |                      |                      |
|               | .note.gnu.property |                      | 64                   |                      |
|               | .note.ABI-tag      |                      | 32                   |                      |
|               | .gnu.hash          |                      | 624                  |                      |
|               | .gnu.version       |                      | 350                  |                      |
|               | .gnu.version_r     |                      | 256                  |                      |
|               | .gcc_except_table  |                      | 694                  |                      |
|               | Total              | 404                  | 2020                 | +500.00%             |
+---------------+--------------------+----------------------+----------------------+----------------------+
| Debug info    | .debug_frame       | 77308                |                      |                      |
|               | .debug_gdb_scripts | 40                   |                      |                      |
|               | .debug_info        | 598195               | 193095               | -67.72%              |
|               | .debug_loc         | 708295               |                      |                      |
|               | .debug_ranges      | 256496               |                      |                      |
|               | .debug_abbrev      | 532                  | 11075                | +2081.77%            |
|               | .debug_line        | 227134               | 49606                | -78.16%              |
|               | .debug_str         |                      | 131399               |                      |
|               | .debug_line_str    |                      | 1620                 |                      |
|               | .debug_rnglists    |                      | 6347                 |                      |
|               | .debug_aranges     |                      | 8896                 |                      |
|               | Total              | 1868000              | 402038               | -78.48%              |
+---------------+--------------------+----------------------+----------------------+----------------------+
| Other         | .strtab            | 46665                | 52863                | +113.28%             |
|               | .shstrtab          | 263                  | 405                  | +153.99%             |
|               | .symtab            | 49224                | 20376                | -58.61%              |
|               | .rela.plt          |                      | 2232                 |                      |
|               | .init_array        |                      | 32                   |                      |
|               | .interp            |                      | 28                   |                      |
|               | .dynsym            |                      | 4200                 |                      |
|               | .dynstr            |                      | 7532                 |                      |
|               | .eh_frame_hdr      |                      | 4132                 |                      |
|               | .fini_array        |                      | 32                   |                      |
|               | .dynamic           |                      | 576                  |                      |
|               | .comment           |                      | 102                  |                      |
|               | .rela.dyn          |                      | 20352                |                      |
|               | .eh_frame          |                      | 18116                |                      |
|               | .preinit_array     |                      | 8                    |                      |
|               | .got               |                      | 48                   |                      |
|               | .data.rel.ro       |                      | 352                  |                      |
|               | .got.plt           |                      | 768                  |                      |
|               | Total              | 96152                | 132154               | +137.44%             |
+---------------+--------------------+----------------------+----------------------+----------------------+
```

## TODO

 * Add styles for tables
 * Add table text alignments
 * Add coloring of the table cells
 * Add table sorting

