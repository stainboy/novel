#!/usr/bin/env bash
set -e


if [[ -z $1 ]]; then
    echo 'Invalid bookId' > /dev/stdout
    exit 1
fi

# load book metadata
source $1.metadata

# workdir
workdir=$bookId
mkdir -p $workdir/tmp || true

# toc
if [ ! -f $workdir/toc.json ]; then
    curl -s "http://m.zongheng.com/h5/ajax/chapter/list?h5=1&bookId=${bookId}&pageNum=1&pageSize=${pageSize}&chapterId=0&asc=0" | jq . > $workdir/toc.json    
fi

# fetch article
toc=($(cat $workdir/toc.json | jq '.chapterlist.chapters[] | .chapterId'))
total="${#toc[@]}"
current=1
for cc in "${toc[@]}"; do
    chapter="$(printf "%04d" $current)__$cc.json"
    if [ ! -f $workdir/tmp/$chapter ]; then
        echo "[${current}/${total}] fetching article..."
        curl -s "http://m.zongheng.com/h5/ajax/chapter?bookId=${bookId}&chapterId=${cc}" -H "Cookie: ___bz=${bz}" > $workdir/tmp/$chapter
    else
        echo "[${current}/${total}] article exists"
    fi
    let current=$current+1
done


# process article
current=1
for cc in "${toc[@]}"; do
    chapter="$(printf "%04d" $current)__$cc.json"
    chapter_md="$(printf "%04d" $current)__$cc.md"
    if [ ! -f $workdir/tmp/$chapter_md ]; then
        echo "[${current}/${total}] process article..."
        chapterName=$(cat $workdir/tmp/$chapter | jq -r .result.chapterName)
        article=$(cat $workdir/tmp/$chapter | jq -r .result.content)
        article="${article//<p>/}"
        article="${article//<\/p>/\\n\\n}"
        echo -e "##${chapterName}##\n\n" > $workdir/tmp/$chapter_md
        echo -e $article >> $workdir/tmp/$chapter_md
    else
        echo "[${current}/${total}] article exists"
    fi
    let current=$current+1
done

# merge
echo -e "#${bookTitle}#\n\n" > "$workdir/${bookTitle}.md"
current=1
for cc in "${toc[@]}"; do
    chapter_md="$(printf "%04d" $current)__$cc.md"
    cat $workdir/tmp/$chapter_md >> "$workdir/${bookTitle}.md"
    let current=$current+1
done

echo 'bye'
