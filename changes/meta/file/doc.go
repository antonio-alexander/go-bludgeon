// Copyright 2022 antonio-alexander. All rights reserved.
// Use of this source code is governed by an MPLv2
// license that can be found in the LICENSE file.

/*
	Package file is a concrete implementation of the Employee interface defined
	in package Meta. File uses an underlying concrete implementation of memory
    meta and utilizes the Serializer functions to write/read the data to disk as
	JSON. File also has a mechanism for file locking that allows a single file
	to be used concurrently by different processes (running in different applications)
	as long as they use the same lock file.
*/
package file
