# df2d(Distribute File to Directory)

## What's df2d??

This df2d command distributes files to a specified location. Here, the meaning of destribution is as follows

1. Copy
2. symbolic link

This is useful for creating a (virtual) environment. What is called CD/CI is useful. However, it's important to remember that the environment is not universal. In other words, Continuous Delivery(CD) is not forever. For example, "The hardware to which software and configuration files are delivered will never change in the same way forever." Is this really true? What about the operating system? Environments change.

However, this is not to say that Continuous Delivery is bad. It's suitable for systems that will be in operation for a long time. It's a powerful way to rebuild an environment, as long as the hardware, software, and other dependencies are the same.

This df2d is aware of continuous delivery, but assumes a situation where the environment changes in a short period of time. For example, this is a scene where a virtual environment(VMware, Virtual Box, docker) is used in the middle of development. This scene is often done on a trial basis and the environment may be reconstructed. At this time, we want to place files easily.

In these scenes, the environment can change dramatically. Do you want to write a Dockerfile in this situation? Automation is not a suitable term for an environment that will no longer be used. This is essentially the same thing as saying that continuous delivery is unsuitable for systems that run for short periods of time.

So, how can we do this easily? It is not efficient to just type `cp` or `ln` commands. The `rsync` command is a good way. Here, df2d is explained together with the `rsync` command. In other words, df2d is positioned similar to the `rsync` command.

The `rsync` command can be used to copy files locally or remotely.

`rsync`

```
                  [remote]
+------------+                 +--------------+
|   server   |                 |    server    |
+------------+                 +--------------+
|   files --------- rsync ---------> files    |
| directories ----- rsync ------> directories |
+------------+                 +--------------+

                   [local]
+---------------------------------------------+
|                   server                    |
+---------------------+-----------------------+
|   files --------- rsync ------->  files     |
| directories ----- rsync ------> directories |
+---------------------+-----------------------+

```

On the other hand, df2d cannot copy to a remote Only in the local environment. The reason why df2d doesn't support remote is that today we often use repository management servers such as Github. In other words, `git clone` is sufficient to copy files between remotes. 

`df2d`

```
                                                     [remote]
   ____________                     +-----------------------------------------+
  (____________)                    |                 server                  |
  |            |                    +-------------------+---------------------+
  |   Github   | --- git clone ---> | files ---------- df2d -------> files    |
  | repository |                    | directories ---- df2d ----> directories |
  |____________|                    +-------------------+---------------------+
        ^
        |
     git push
        |
  +-------------+
  |    files    | [local]
  | directories |
  +-------------+

```

There are other differences between df2d and `rsync`. That is how to specify the destination and source of the copy. The `rsync` command specifies the destination and source of the copy with arguments. df2d is described in the toml file. If only a small number of files or directories are to be copied, it is not difficult to specify them as arguments. However, there will be more as development progresses, including executables, configuration files. For example, this is a situation that is likely to occur in implementing microservices in Go.

The number of files can be reduced by managing them not by file units, but by dividing them by directory units. However, the more systems and applications that are linked together other than microservices, the more there will be. When trying to use the `rsync` command in such cases, there is a way to create a shell script to automate the input of the destination and source of the copy. However, this has a drawback: we have to check the shell script to see which directories (files) to put where. This is not good. 

For each directory, df2d is described in a toml file where to place it. In other words, by looking at the directory, we can tell at a glance where it will be placed. Also, there is no need to create even a shell script, just write a few lines of a configuration file and place it in the directory from which it will be copied.

```
     d2f2
       | _______               _______
       | | src |               | dst |
 +-----v-+-----+       +-------------+
 |  df2d.toml  |       |             |
 |    file1 -------------> file1     |
 |    file2 -------------> file2     |
 |     ...     |       |    ...      |
 |    dir1  -------------> dir1      |
 |     ...     |       |    ...      |
 +-------------+       +-------------+
```

df2d is a tool that copies on a per-directory basis. Therefore, df2d allows developers to proceed with development while organizing the location of files. df2d is, at worst, a degraded version of the rsync command. Not everything is as it seems. The table below shows how df2d should be used, including what has been explained so far.

tools|development period|
:--|:--|
cp,mv|short|
df2d|middle|
rsync|long|
CD(Continuous Development)|...|

It is important to use them differently.

df2d doesn't have the function to delete. Finally, the reason for this is explained. Copying and deleting files are opposite properties. Deleting is actually not easy. If a necessary file is deleted without permission, the developer will be in trouble. The `rsync` command controls what to delete with the `-exclude` option, for example. With this option, the `rsync` command can act as if the source and destination are synchronized. This means that the `rsync` command will not delete a file if there is no delete list, even if it is an `rsync` command. That is how careful deletion should be.

This is not to say that df2d doesn't need a delete function. As to why df2d doesn't implement the delete function, from the table above, df2d is in a rougher position to use than rsync (middle). In other words, df2d states that it is better to leave unused files in place without deleting them. Of course, as development progresses, unnecessary files must be deleted. However, developers need not be nervous in the early stages of development. Also, although there is a possibility that unnecessary files will be mass-produced at the copy destination, if unnecessary files are removed from the repository, only necessary files will exist. In other words, the developer only needs to look at the repository to see which files are necessary. This is why df2d doesn't have a delete function.

## How2Use

1. `go build .`
2. `./df2d -e env.toml`

### Configuration file

There are two main types of df2d configuration files.

1. env.toml
2. df2d.toml

Both are toml files and their uses are as follows:

config|summary
:--|:--
env.toml|Defines what is handled in the source directory. It is possible to specify files that are not to be copied.
df2d.toml|Define the specific copy destination. It is possible to determine the means of copying.

#### env.toml

Let df2d load with the -e option. The parameters of env.toml are as follows.

- **root**
  - The root directory of the copy source.
- **excluded_ext**
  - Extended files included in this list are not copied.
- **isExcluded_hidden_file**
  - `true`:Prefixes with `.` will not be copied.

#### df2d.toml

The file name must be `df2d`. Put it in the directory from which it was copied.

- **isMkdir**
  - `true`:If the directory does not exist at the copy destination, create it.
- **dest_dir**
  - Specify the destination directory.
- **isOverwrite**
  - `true`:Overwrites an already existing file.
- **isSymLink**
  - `true`:Create a symbolic link, not a copy. Symbolic links to the directories will not be made.

### Logic

The logic for reading df2d.toml is as follows:

1. Copy the files and directories in the directory where df2f.toml is located.
2. The df2d.toml in the lowest directory of the directory tree takes precedence.

<details>
<summary>directory tree</summary>

```
df2d-src
|-- README.md
|-- .ignore
|   \`-- IGNORE.md
|-- depth1-1
|   |-- README2-1.md
|   \`-- depth2-1
|       \`-- README3-1.md
|-- depth1-2
|   |-- .ignore
|   |-- README2-2.md
|   \`-- df2d.toml
\`-- df2d.toml

     |
     | df2d
     v

df2d-dst
|-- README.md
|-- depth1-1
|   |-- README2-1.md
|   `-- depth2-1
|       `-- README3-1.md
`-- new
    `-- README2-2.md
```
</details>

The configuration file for this sample is as follows:

`env.toml`

```
root="/your/path/to/df2d-src"
isExcluded_hidden_file=true
```

`df2d-src/df2d.toml`

```
isMkdir=true
dest_dir="/your/path/to/df2d-dst/"
isOverwrite=true
```

`df2d-src/depth1-2/df2d.toml`

```
isMkdir=true
dest_dir="/your/path/to/df2d-dst/new/"
isOverwrite=true
```

