import React, { useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import i18n, { packageNS } from '../../i18n';

const MAX_FILE_SIZE_MEGABYTES = 10;

function UserProfilePicFile(props) {
  const onDrop = useCallback((files) => {
    const { onChange } = props;
    let output = {};

    files.forEach((file) => {
      const reader = new FileReader();

      reader.onabort = (err) => {
        output.errorMessage = i18n.t(`${packageNS}:tr000435`, { error: err });
        onChange && onChange(output);
      };

      reader.onerror = (err) => {
        output.errorMessage = i18n.t(`${packageNS}:tr000436`, { error: err });
        onChange && onChange(output);
      };

      reader.onload = () => {
        output.result = reader.result;
        output.successMessage = i18n.t(`${packageNS}:tr000437`, { filename: file.name });
        onChange && onChange(output);
      };

      // base64 is format of readAsDataURL
      reader.readAsDataURL(file);
    });
  }, []);

  const onDropRejected = useCallback((files) => {
    const { onChange } = props;
    let output = {};
    const errorMessage = i18n.t(`${packageNS}:tr000438`, { maxFileSizeMegabytes: MAX_FILE_SIZE_MEGABYTES, filename: files[0].name });//`Maximum file upload size of ${MAX_FILE_SIZE_MEGABYTES} MB exceeded by file ${files[0].name}`;
    console.error(errorMessage);
    output.errorMessage = errorMessage;
    onChange && onChange(output);
  }, []);

  const { getRootProps, getInputProps } =
    useDropzone({
      // Accept all image formats
      accept: 'image/*',
      // Max image size that may be uploaded in bytes
      maxSize: MAX_FILE_SIZE_MEGABYTES * 1000000,
      onDrop,
      onDropRejected
    });

  return (
    <div className='input-file' style={{ padding: "10px" }} {...getRootProps()}>
      <input {...getInputProps()} />
      <span>{props.profilePicImage}</span>
      <span className='label' style={{ marginLeft: "20px", marginRight: "10px" }}>{i18n.t(`${packageNS}:tr000434`)}</span>
      <i className="mdi mdi-upload"></i>
    </div>
  );
}

export default UserProfilePicFile;
