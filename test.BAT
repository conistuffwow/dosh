SET COLOR=green
SET LANG=Go

ECHO The Color is %COLOR%
ECHO The Language is %LANG%
ECHO The Color is %COLOR% and the Language is %LANG%

IF %LANG%==Go ECHO This should be printed
IF %COLOR%==red ECHO This should not be printed

ECHO jumping to label a
GOTO a

:a

ECHO This is label a
ECHO The Color is still %COLOR%
ECHO The Language is still %LANG%