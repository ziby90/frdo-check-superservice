package error_handler

var ValidBase64Err = ErrorType{Type: 1, ToUserType: 500}
var ValidJsonErr = ErrorType{Type: 2, ToUserType: 500}
var ValidDataErr = ErrorType{Type: 3, ToUserType: 407}
var ValidXmlErr = ErrorType{Type: 4, ToUserType: 404}
var BDErr = ErrorType{Type: 5, ToUserType: 500}
var BDSeqErr = ErrorType{Type: 5, ToUserType: 500}
var BDNotFound = ErrorType{Type: 6, ToUserType: 404}
var BDFoundExist = ErrorType{Type: 10, ToUserType: 408}
var BDFoundRelation = ErrorType{Type: 7, ToUserType: 405}
var MergeErr = ErrorType{Type: 8, ToUserType: 500}
var CopyErr = ErrorType{Type: 9, ToUserType: 500}
