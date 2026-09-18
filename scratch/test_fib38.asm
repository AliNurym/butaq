; Butaq compiled output — NASM x86-64
; Generated automatically. Do not edit.

    default rel

section .data
    fmt_float db "%g", 10, 0
    fmt_str   db "%s", 10, 0
    fmt_int   db "%lld", 10, 0
    fmt_numstr db "%g", 0
    fmt_input  db "%[^\n]", 0
    file_r    db "r", 0
    file_w    db "w", 0
    newline   db 10, 0
    err_null  db "Қате: бос сілтемеге (null) жүгіну!", 0
    err_div0  db "Қате: нөлге бөлуге болмайды!", 0
    err_bounds db "Қате: жиын/тізім шекарасынан шығу!", 0
    str_true_val  db "ақиқат", 0
    str_false_val db "жалған", 0
    flt_1 dq 0x3FF0000000000000  ; 1
    flt_2 dq 0x3FF0000000000000  ; 1
    flt_3 dq 0x4000000000000000  ; 2
    flt_4 dq 0x4043000000000000  ; 38
    flt_5 dq 0x408F400000000000  ; 1000
    str_1 db "RESULT:", 0
    str_2 db "TIME_MS:", 0
section .bss
    scratch_buf resb 65536
    input_buf   resb 4096
    numstr_buf  resb 64
    char_buf    resb 4
    global_var_бастау resq 1
    global_var_нәтиже resq 1
    global_var_аяқтау resq 1
    global_var_мс resq 1

section .text
    global main
    extern printf
    extern strlen
    extern strcmp
    extern strcat
    extern strcpy
    extern malloc
    extern free
    extern calloc
    extern strdup
    extern sprintf
    extern fopen
    extern fclose
    extern fread
    extern fwrite
    extern fseek
    extern ftell
    extern fflush
    extern scanf
    extern atof
    extern rand
    extern SetConsoleOutputCP
    extern _runtime_clock
    extern ExitProcess
    extern _thread_spawn
    extern builtin_init_args
    extern builtin_result_ok_num
    extern builtin_result_ok_ptr
    extern builtin_result_err
    extern builtin_result_is_error
    extern builtin_result_get_num
    extern builtin_result_get_ptr
    extern builtin_result_get_err
    extern builtin_cos
    extern builtin_num_to_str
    extern builtin_mutex_destroy
    extern builtin_string_builder_destroy
    extern builtin_json_get_object
    extern builtin_mutex_lock
    extern builtin_string_builder_create
    extern builtin_str_split_count
    extern builtin_math_atan2
    extern builtin_base64_encode
    extern builtin_math_log2
    extern builtin_map_size
    extern builtin_math_round
    extern builtin_args_count
    extern builtin_json_list_size
    extern builtin_math_log
    extern builtin_sqrt
    extern builtin_math_exp
    extern builtin_file_exists
    extern builtin_json_get_string
    extern builtin_json_add_object
    extern builtin_json_add_list
    extern builtin_string_builder_append_number
    extern builtin_str_starts_with
    extern builtin_str_to_num
    extern builtin_time_seconds
    extern builtin_input
    extern builtin_map_has_key
    extern builtin_mutex_create
    extern builtin_str_lower
    extern builtin_str_ends_with
    extern builtin_math_log10
    extern builtin_math_tan
    extern builtin_pow
    extern builtin_md5
    extern builtin_json_add_string
    extern builtin_json_add_bool
    extern builtin_file_delete
    extern builtin_json_get_list
    extern builtin_str_trim
    extern builtin_math_min
    extern builtin_json_write
    extern builtin_json_list_get
    extern builtin_map_get_string
    extern builtin_string_builder_append_string
    extern builtin_rand_float
    extern builtin_exec_output
    extern builtin_sha256
    extern builtin_map_delete_key
    extern builtin_math_abs
    extern builtin_sin
    extern builtin_math_pi
    extern builtin_system
    extern builtin_json_parse
    extern builtin_thread_join
    extern builtin_string_builder_append_char
    extern builtin_str_replace
    extern builtin_math_max
    extern builtin_sleep_seconds
    extern builtin_json_create
    extern builtin_json_add_number
    extern builtin_map_get_number
    extern builtin_string_builder_to_string
    extern builtin_str_upper
    extern builtin_str_index_of
    extern builtin_math_floor
    extern builtin_base64_decode
    extern builtin_json_get_bool
    extern builtin_mutex_unlock
    extern builtin_str_split_get
    extern builtin_utf8_str_len
    extern builtin_json_get_number
    extern builtin_map_set_string
    extern builtin_str_slice
    extern builtin_arg_get
    extern builtin_utf8_char_at
    extern builtin_time_str
    extern builtin_map_create
    extern builtin_map_set_number
    extern builtin_math_ceil
    extern builtin_math_atan

main:
    push rbp
    mov rbp, rsp
    sub rsp, 32
    sub rsp, 32
    call builtin_init_args
    add rsp, 32
    ; бастау болсын
    ; уақыт
    sub rsp, 32
    call _runtime_clock
    add rsp, 32
    movsd [global_var_бастау], xmm0
    ; нәтиже болсын
    ; шақыру фибо
    movsd xmm0, [flt_4]
    movq rcx, xmm0
    sub rsp, 32
    call фибо
    add rsp, 32
    movsd [global_var_нәтиже], xmm0
    ; аяқтау болсын
    ; уақыт
    sub rsp, 32
    call _runtime_clock
    add rsp, 32
    movsd [global_var_аяқтау], xmm0
    ; мс болсын
    movsd xmm0, [global_var_аяқтау]
    mov rax, [global_var_аяқтау]
    subsd xmm0, [global_var_бастау]
    mulsd xmm0, [flt_5]
    movsd [global_var_мс], xmm0
    ; жазу
    lea rax, [str_1]
    mov rdx, rax
    lea rcx, [fmt_str]
    sub rsp, 32
    call printf
    add rsp, 32
    xor rcx, rcx
    sub rsp, 32
    call fflush
    add rsp, 32
    ; жазу
    movsd xmm1, [global_var_нәтиже]
    movq rdx, xmm1
    lea rcx, [fmt_float]
    sub rsp, 32
    call printf
    add rsp, 32
    xor rcx, rcx
    sub rsp, 32
    call fflush
    add rsp, 32
    ; жазу
    lea rax, [str_2]
    mov rdx, rax
    lea rcx, [fmt_str]
    sub rsp, 32
    call printf
    add rsp, 32
    xor rcx, rcx
    sub rsp, 32
    call fflush
    add rsp, 32
    ; жазу
    movsd xmm1, [global_var_мс]
    movq rdx, xmm1
    lea rcx, [fmt_float]
    sub rsp, 32
    call printf
    add rsp, 32
    xor rcx, rcx
    sub rsp, 32
    call fflush
    add rsp, 32
    ; --- exit ---
    xor eax, eax
    add rsp, 32
    pop rbp
    ret

; функция фибо
фибо:
    push rbp
    mov rbp, rsp
    sub rsp, 16
    movsd [rbp-8], xmm0  ; param num n
    ; егер
    movsd xmm0, [rbp-8]
    mov rax, [rbp-8]
    comisd xmm0, [flt_1]
    ja .else_1
    movsd xmm0, [rbp-8]
    mov rax, [rbp-8]
    ; қайтару
    mov rsp, rbp
    pop rbp
    ret
.else_1:
    ; шақыру фибо
    movsd xmm0, [rbp-8]
    mov rax, [rbp-8]
    subsd xmm0, [flt_2]
    movq rcx, xmm0
    sub rsp, 32
    call фибо
    add rsp, 32
    sub rsp, 8
    movsd [rsp], xmm0
    ; шақыру фибо
    movsd xmm0, [rbp-8]
    mov rax, [rbp-8]
    subsd xmm0, [flt_3]
    movq rcx, xmm0
    sub rsp, 32
    call фибо
    add rsp, 32
    movsd xmm1, [rsp]
    add rsp, 8
    addsd xmm1, xmm0
    movsd xmm0, xmm1
    ; қайтару
    mov rsp, rbp
    pop rbp
    ret
    mov rsp, rbp
    pop rbp
    ret

; --- Runtime error handlers ---
_runtime_null_pointer_error:
    lea rdx, [err_null]
    lea rcx, [fmt_str]
    sub rsp, 32
    call printf
    add rsp, 32
    mov rcx, 1
    call ExitProcess
_runtime_divide_by_zero_error:
    lea rdx, [err_div0]
    lea rcx, [fmt_str]
    sub rsp, 32
    call printf
    add rsp, 32
    mov rcx, 1
    call ExitProcess
_runtime_array_bounds_error:
    lea rdx, [err_bounds]
    lea rcx, [fmt_str]
    sub rsp, 32
    call printf
    add rsp, 32
    mov rcx, 1
    call ExitProcess
