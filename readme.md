# mrglass

small util program to correlate hashstack output with usernames

```
Usage: 
	mrglass <USERHASH> <CRACKED>

	<CRACKED> | mrglass <USERHASH> 

	USERHASH is a newline-seperated file in the the format of 
	'USER:HASH'

	CRACKED is a newline-seperated file or pipe whose entries 
	are in the the format of 
	'HASH:PLAINTEXT_PASSWORD'

Examples:
	mrglass hashes_with_usernames.txt cracked.txt
	hashstack lists cracked 1 1 | mrglass hashes_with_usernames.txt
	mrglass hashes_with_usernames.txt <(hashstack lists cracked $pID $lID)
```
