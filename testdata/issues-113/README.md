This directory was created to try to bring to the fore the issue outlined in; https://github.com/komuw/meli/issues/113  


Basically, meli seems to be allocating a lot of memory. And the bulk of that is when during image building, meli 
tars every file inside docker context.   
So to my guess is that the more files we have in the context the easier for that problem to manifest.  
Hence we create this folder with a lot of files; and the more nested they are , the better. 
