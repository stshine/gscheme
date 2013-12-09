(define insertL
  (lambda (new old lat)
    (cond
     ((null? lat) lat)
     ((eq? old (car lat)) (cons new lat))
     (else (cons  (car lat) (insertL new old (cdr lat)))))))

(insertL (quote sev) (quote grade) (quote (loop grade)))
