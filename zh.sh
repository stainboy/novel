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
curl -s "http://m.zongheng.com/h5/ajax/chapter/list?h5=1&bookId=${bookId}&pageNum=1&pageSize=${pageSize}&chapterId=0&asc=0" > $workdir/toc.json

# fetch article
toc=($(cat $workdir/toc.json | jq '.chapterlist.chapters[] | .chapterId'))
total="${#toc[@]}"
current=1
for cc in "${toc[@]}"; do
    chapter="$(printf "%04d" $current)__${cc}_1.json"
    if [ ! -f $workdir/tmp/$chapter ]; then
        echo "[${current}/${total} (1/?)] fetching article..."
        curl -s "http://m.zongheng.com/h5/ajax/chapter?bookId=${bookId}&chapterId=${cc}" -H "Cookie: ___bz=${bz}" > $workdir/tmp/$chapter
    else
        echo "[${current}/${total} (1/?)] article exists"
    fi

    # paging >_<!
    pageCount=$(cat $workdir/tmp/$chapter | jq .result.pageCount)
    for (( i = 2; i < $pageCount+1; i++ )); do
        chapter_x="$(printf "%04d" $current)__${cc}_$i.json"
        if [ ! -f $workdir/tmp/$chapter_x ]; then
            echo "[${current}/${total} ($i/$pageCount)] fetching article..."
            curl -s "http://m.zongheng.com/h5/ajax/chapter?bookId=${bookId}&chapterId=${cc}_${i}" -H "Cookie: ___bz=${bz}" > $workdir/tmp/$chapter_x
        else
            echo "[${current}/${total} ($i/$pageCount)] article exists"
        fi
    done

    let current=$current+1
done


# process article
current=1
for cc in "${toc[@]}"; do
    chapter="$(printf "%04d" $current)__${cc}_1.json"
    chapter_md="$(printf "%04d" $current)__$cc.md"

    echo "[${current}/${total}] process article..."
    chapterName=$(cat $workdir/tmp/$chapter | jq -r .result.chapterName)
    echo -e "##${chapterName}##\n\n" > $workdir/tmp/$chapter_md

    pageCount=$(cat $workdir/tmp/$chapter | jq .result.pageCount)
    for (( i = 1; i < $pageCount+1; i++ )); do
        chapter_x="$(printf "%04d" $current)__${cc}_$i.json"
        article=$(cat $workdir/tmp/$chapter_x | jq -r .result.content)
        article="${article//<p>/}"
        article="${article//<\/p>/\\n\\n}"
        echo -e $article >> $workdir/tmp/$chapter_md
        echo -e '\n\n' >> $workdir/tmp/$chapter_md
    done

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
