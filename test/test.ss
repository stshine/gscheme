(define factor
  (lambda (x)
    (cond
     ((zero? (- x 1)) 1)
     (else (* x (factor (- x 1)))))))

(factor 4)
