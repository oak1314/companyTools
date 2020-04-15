import os
from comtypes.client import CreateObject


class ComTypes:
    def ppt_pdf(self,path):
        # PPT 转 PDF
        pdf_path = path.replace('ppt', 'pdf') # pdf保存路径 推荐使用绝对路径
        try:
            p = CreateObject("PowerPoint.Application")
            ppt = p.Presentations.Open(path)
            ppt.ExportAsFixedFormat(pdf_path, 2, PrintRange=None)
            ppt.Close()
            p.Quit()
        except Exception as e:
            pass
    

    def word_pdf(self,path):
        # Word转pdf
        pdf_path = path.replace('doc', 'pdf')
        w = CreateObject("Word.Application")
        doc = w.Documents.Open(path)
        doc.ExportAsFixedFormat(pdf_path, 17)
        doc.Close()
        w.Quit()


    def excel_pdf(self,path):
        # Excel转pdf
        pdf_path = path.replace('xlsx', 'pdf')
        xlApp = CreateObject("Excel.Application")
        book = xlApp.Workbooks.Open(path)
        book.ExportAsFixedFormat(0, pdf_path)
        book.Close(False)
        xlApp.Quit()
        
    def file_name(self,file_dir): 
        L=[]
        for entry in os.listdir(file_dir):
            # 使用os.path.isfile判断该路径是否是文件类型
            if os.path.isfile(os.path.join(file_dir, entry)):
                if os.path.splitext(entry)[1] == '.xlsx':
                    L.append(os.path.join(file_dir, entry))
                # print(entry)

        # for root, dirs, files in os.walk(file_dir):
        #     for file in files:
        #         if os.path.splitext(file)[1] == '.xlsx':
        #             L.append(os.path.join(root, file))
        return L

if __name__ == '__main__':
    comty  = ComTypes()
    for xlsx_file in comty.file_name(os.getcwd()):
        comty.excel_pdf(xlsx_file)

