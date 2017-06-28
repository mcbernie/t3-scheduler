# t3-scheduler
In „Dateiliste“ liegt dann in Zukunft eine Datei namens scheduler.txt
Die Datei sieht ungefähr aus:
[wacht]
pages=132,33,6,8,10

[receivers]
132=nico.brüggemann,katrin.uhrbrock
33=nico.brüggemann
...

[mail]
header=Content in page %pagename% updated
body=Content in %pagename% updated

Diese Datei kannst du dann frei anpassen.
Bei pages kommen alle ids der pages rein die „überwacht“ werden sollen.
In receivers links vor dem = die page id und rechts die benutzernamen der benutzer die informiert werden sollen wenn der inhalt der angegebenen seite geändert wurde.
Im Abschitt mail kannst du den betreff unter header anpassen und den inhalt kannst du in body anpassen.
In bedien werten hast du den platzhalter %pagename% für den seitennamen und %user% für den Benutzernamen zur verfügung.
