import codecs
import sys
import json
import tqdm
import os

# from elasticsearch import Elasticsearch
# from elasticsearch.helpers import bulk

data_dir = "H:\\Scholar"
paper_path = os.path.join(data_dir)

bad_data_num =0

file_list = []
author_list = []
author_dict = {}

venue_list = []
venues ={}
journal_list = []
journals = {}
journal_out_list = []

in_citation_dict = {}
in_citation_list = []

for paper_file in os.listdir(paper_path):
    domain = os.path.abspath(paper_path)
    if paper_file.startswith("s2-corpus") and  not paper_file.endswith(".gz"):
        paper_file = os.path.join(domain,paper_file)
        file_list.append(paper_file)
file_list.sort()

# file_list = ["H:\\Scholar\\s2-corpus-000"]
def proc_file_list(file_list,label):
    global  bad_data_num
    file_num = 0
    for file in file_list:
        with codecs.open(file,'r','utf-8') as f:

            i = 0
            for line in f.readlines():
                if len(line.strip()) == 0:
                    break

                try:
                    dict_item = json.loads(line)
                except:
                    bad_data_num += 1
                    continue
                list_itme = list(dict_item["fieldsOfStudy"])
                if "Computer Science" not in list_itme :
                    continue
                if label == "authors":
                    for author in dict_item['authors']:
                        # print(author)
                        if (str(author['ids']) in author_dict.keys()) == True:
                            author_dict[str(author['ids'])]["citation_num"] += len(dict_item['inCitations'])
                            author_dict[str(author['ids'])]["publish_num"] += 1

                        else:
                            item = {"author_id" : list(author['ids']),"publish_num" : 1,"name":author['name'],'org':"","citation_num":len(dict_item['inCitations'])}
                            author_dict[str(author['ids'])] = item
                elif label == "journal":
                    venue,journalName,journalVolume,journalPages = dict_item['venue'],dict_item['journalName'],dict_item['journalVolume'],dict_item['journalPages']
                    journalName = str(journalName).strip()
                    if journalName != "":
                        if journalName not in journal_list:
                            item = {"name":journalName,"volume":journalVolume,"pages":journalPages,'id':len(journal_list),'venue':venue}
                            item['publish_num'],item['citation_num'] = 1,len(dict_item['inCitations'])

                            item['authors'] = []
                            for author in dict_item['authors']:
                                item['authors'].append(str(author['ids']))
                            item['authors_num'] = len(item['authors'])  #
                            journal_list.append(journalName)
                            journals[journalName] = item
                        else:
                            journals[journalName]['publish_num'], journals[journalName]['citation_num'] = journals[journalName]['publish_num']+1, journals[journalName]['citation_num']+len(dict_item['inCitations'])
                            for author in dict_item['authors']:
                                if (str(author['ids']) not in journals[journalName]["authors"] ) :
                                    journals[journalName]['authors'].append(str(author['ids']))
                            journals[journalName]['authors_num'] = len(journals[journalName]['authors'])  # 投稿人数

                elif label == "inCitations":
                    item = {"id":dict_item["id"],"inCitations":dict_item["inCitations"],"inCitationsNum":len(dict_item["inCitations"])}
                    # in_citation_dict[dict_item["id"]] = item
                    in_citation_list.append(item)
                i += 1
            file_num += 1


        if label == "inCitations":
            print(file_num,len(in_citation_dict.keys()))
            write_inCitations_list()
        elif label == "authors":
            print(file_num,len(author_dict))
        elif label == "journal":
            print(file_num, len(journal_list))


def make_author_list():
    print(len(list(author_dict.keys())))
    for key in list(author_dict.keys()):
        author_list.append(author_dict.pop(key))
def write_author_list():
    with codecs.open(data_dir+"authors.txt","w") as f:
        for author in author_list:
            f.write(json.dumps(author)+"\n")
    author_list.clear()
def make_journal_out_list():
    print(len(list(journals.keys())))
    for key in list(journals.keys()):
        journal_out_list.append(journals.pop(key))

def write_journal_out_list():
    with codecs.open(data_dir+"journal.txt","w") as f:
        for journal in journal_out_list:
            del journal['authors']
            f.write(json.dumps(journal)+"\n")
    journal_out_list.clear()
def write_inCitations_list():
    with codecs.open(data_dir+"inCitations.txt","a") as f:
        for item in in_citation_list:
            f.write(json.dumps(item)+"\n")
    in_citation_list.clear()
proc_file_list(file_list,"authors")
make_author_list()
write_author_list()

proc_file_list(file_list,"journal")
make_journal_out_list()
write_journal_out_list()

proc_file_list(file_list,"inCitations")