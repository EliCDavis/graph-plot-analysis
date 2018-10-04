# Graph Plot Anaylsis

Examine points on a graph and determine their position and standard deviation on the y-axis. Util built for my sister. Output graph is provided so a human can check if all the points got highlighted.

Tested with Golang 1.10

To get working:
```bash
git clone https://github.com/EliCDavis/graph-plot-analysis.git
cd graph-plot-analysis
go get github.com/fogleman/gg
go build graphanal.go
```

If you want to use the tool anywhere in your terminal, be sure to add the folder location to your environment variables path or bash_rc.

```
Usage of graphanal:
  -in string
        Name of the image file to examine (default "input.png")
  -out string
        Name of the image output for human checking (default "out.png")
  -threshold int
        Adjust this value if some of your points arn't being found (default 24000)
```

If you have a file named `input.png` in the directory you're working out of you can just run `graphanal`.

Example of more complicated input:
```
graphanal -in 001.png -threshold 15000
```

# Example Input

input.png

![input](https://i.imgur.com/7fZY1MG.png)

# Example Output
```
2018/10/04 14:51:57 Points found: 228
2018/10/04 14:51:57 Standard Deviation: 0.098147
```

![output](https://i.imgur.com/KqEDoFI.png)