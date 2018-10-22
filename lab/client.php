<?php

$res = fsockopen("localhost",1337);
$message = "Hello World";
fwrite($res,package($message));
$x = true;
$ret   = "";
$check = false;
while($x == true) {
    $tmp = fread($res,1);
    if(strlen($tmp) > 0) {
        $ret .= $tmp;
    }
    if($check == false  && strlen($ret) > 3) {
        $packageLength = 403 - intval(substr($ret,0,3));
        $check = true;
    }
    if($check == true && strlen($ret) == $packageLength + 3) {
        var_dump($ret);
        break;
    }
}

function package($str) {
    return (403-strlen($str)) . $str;
}





?>