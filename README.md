# sim-pad-data-tools

Tools for sorting and combining the data of a resusci anne training puppet

### Speciality for MacOS
- If you are running on a new **M1** Mac use the binary called 
sim-pad-data-tools-m1

- If you are running on an old **Intel** Mac use the binary called
sim-pad-data-tools-darwin-amd64


### Log override
In case you want to override the Log create a custom json file named `log-override.json`
in the same directory as the `EventLog.xml` file
<br>
Fill the file with the following data:
```json
{
  "scenario": "",
  "instructor": "",
  "group": "",
  "case": ""
}
```

If you only want to change one value (or more) you only need to add the entry you want to change:
```json
{
  "scenario": "Lungenembolie",
  "group": "C1"
}
```
Empty strings `""` will be ignored by the tool. <br>
Please keep in mind that in the json format the last comma must be omitted. <br>
Otherwise it's not valid json and will produce an error.
