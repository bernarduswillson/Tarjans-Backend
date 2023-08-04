# LZW + RLE compression and decompression🎲
>Tugas Seleksi IRK
## Table of Contents
* [Contributors](#contributors)
* [General Information](#general-information)
* [Local Setup](#local-setup)
* [File Input Structure](#file-input-structure)
## Contributors
| NIM | Nama |
| :---: | :---: |
| 13521021 | Bernardus Willson  |
## General Information 
LZW compression is a lossless data compression algorithm, meaning that no data is lost during the compression and decompression process. It works by creating a table of all the unique strings that appear in the input data. When a new string is encountered, it is looked up in the table. If the string is found in the table, its corresponding code is emitted. If the string is not found in the table, it is added to the table and a new code is emitted.

RLE stands for Run-Length Encoding. It is a lossless data compression algorithm that works by replacing runs of consecutive data values with a single value and a count. For example, if a sequence of data contains 100 consecutive zeros, then the RLE algorithm would store this as a single value of 0 and a count of 100.

This project is a web-based application that can compress and decompress text by combining LZW and RLE algorithms. Original text is compressed using the LZW algorithm and then it generates 8 bit binary numbers. The numbers are then compressed again using RLE.
## Local Setup
<br>
1. Clone FE and BE repo using the command below: 

```
git clone https://github.com/bernarduswillson/LZW-Frontend
```
```
git clone https://github.com/bernarduswillson/LZW-Backend
```
<br>
2. Install dependencies :

```
yarn
```
<br>
3. Run localhost server :
<br>
for FE:

```
yarn dev
```
for BE:

```
npm run app
```
Alternatively, you can open the website by using this link
```
https://lzw-frontend.vercel.app/
```

![](doc/home.png)
![](doc/home2.png)

<br>
4. Type the text you want to compress or decompress in the text area and click the compress button. The compressed or decompressed text will be shown in the text area below.

![](doc/input.png)

<br>
5. You can also save the compressed or decompressed text to a file by clicking the save button. The history of the compressed or decompressed text will be shown in the history section on the bottom of the page.